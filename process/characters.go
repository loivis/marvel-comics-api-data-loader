package process

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/avast/retry-go"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
	"github.com/loivis/marvel-comics-api-data-loader/marvel/mclient/operations"
	"github.com/loivis/marvel-comics-api-data-loader/marvel/models"
)

func (p *Processor) loadCharacters(ctx context.Context) error {
	var err error

	err = p.loadAllCharactersWithBasicInfo(ctx)
	if err != nil {
		return fmt.Errorf("error loading all characters info: %v", err)
	}

	log.Info().Msg("all characters loaded")

	err = p.complementAllCharacters(ctx)
	if err != nil {
		return fmt.Errorf("error complementing all characters: %v", err)
	}

	log.Info().Msg("all characters complemented")

	return nil
}

func (p *Processor) loadAllCharactersWithBasicInfo(ctx context.Context) error {
	remote, err := p.getCharacterCount(ctx)
	if err != nil {
		return fmt.Errorf("error fetching character count: %v", err)
	}
	log.Info().Str("type", "character").Int("count", remote).Msg("character count from api")

	existing, err := p.store.GetCount(ctx, "characters")
	if err != nil {
		return err
	}
	log.Info().Str("type", "character").Int("count", existing).Msg("existing character count")

	if int(remote) == existing {
		log.Info().Int("local", existing).Int("remote", remote).Msg("no missing characters")
		return nil
	}

	log.Info().Int("local", existing).Int("remote", remote).Msg("missing characters, reload")

	return p.loadMissingCharacters(ctx, int32(existing), int32(remote))
}

func (p *Processor) getCharacterCount(ctx context.Context) (int, error) {
	var limit int32 = 1
	params := &operations.GetCharactersCollectionParams{
		Limit: &limit,
	}
	p.setParams(ctx, params)

	col, err := p.mclient.Operations.GetCharactersCollection(params)
	if err != nil {
		return 0, err
	}

	return int(col.Payload.Data.Total), nil
}

func (p *Processor) loadMissingCharacters(ctx context.Context, starting, count int32) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	charCh := make(chan *m27r.Character, int32(p.concurrency)*p.limit)
	conCh := make(chan struct{}, p.concurrency)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	go func() {
		var characters []*m27r.Character
		defer func() {
			doneCh <- struct{}{}
		}()

		batchSave := func(characters []*m27r.Character) error {
			err := retry.Do(func() error {
				return p.store.SaveCharacters(ctx, characters)
			})

			if err != nil {
				return err
			}

			log.Info().Int("count", len(characters)).Msg("batch saved characters")

			return nil
		}

		for character := range charCh {
			characters = append(characters, character)

			if len(characters) >= p.storeBatch {
				if err := batchSave(characters); err != nil {
					errCh <- err
					break
				}
				characters = []*m27r.Character{}
			}
		}

		batchSave(characters)
	}()

	var g errgroup.Group

	for i := int(starting / p.limit); i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			select {
			case err := <-errCh: // check if any error saving data
				cancel()
				log.Info().Int32("offset", offset).Msgf("cancelled fetching paged characters: %v", err)
				return fmt.Errorf("cancelled fetching paged characters limit %d offset %d: %v", p.limit, offset, err)
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("offset", offset).Msg("ctx already cancelled")
				return nil
			default: // default to avoid blocking
			}

			params := &operations.GetCharactersCollectionParams{
				Limit:  &p.limit,
				Offset: &offset,
			}
			p.setParams(ctx, params)

			err := retry.Do(
				func() error {
					col, err := p.mclient.Operations.GetCharactersCollection(params)
					if err != nil {
						return err
					}

					for _, res := range col.Payload.Data.Results {
						character, err := convertCharacter(res)
						if err != nil {
							return err
						}

						charCh <- character
					}

					log.Info().Int32("offset", offset).Int32("count", col.Payload.Data.Count).Msg("fetched paged characters")

					return nil
				},
				retry.OnRetry(retryLog(offset)),
				retry.RetryIf(retryIf(offset)),
			)

			if err != nil {
				cancel()
				log.Info().Int32("offset", offset).Msg("cancelled fetching paged characters")
				return fmt.Errorf("error fetching with limit %d offset %d: (%T) %v", p.limit, offset, err, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	close(charCh)

	select {
	case <-doneCh:
		log.Info().Msg("fetched all missing characters with basic info")
	}

	return nil
}

func (p *Processor) complementAllCharacters(ctx context.Context) error {
	ids, err := p.store.IncompleteIDs(ctx, "characters")
	if err != nil {
		return fmt.Errorf("error get imcomplete character ids: %v", err)
	}

	if len(ids) == 0 {
		log.Info().Msg("no incomplete character")
		return nil
	}

	log.Info().Int("count", len(ids)).Msg("fetched incomplete character ids")

	var g errgroup.Group

	conCh := make(chan struct{}, p.concurrency)

	for _, id := range ids {
		conCh <- struct{}{}

		id := id

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			character, err := p.getCharacterWithFullInfo(ctx, id)
			if err != nil {
				return fmt.Errorf("error fetching character %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("fetched character with full info converted")

			err = p.store.SaveOne(ctx, character)
			if err != nil {
				return fmt.Errorf("error saving character %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("saved character")

			log.Info().Int("id", id).Msgf("complemented character")

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error complementing characters: %v", err)
	}

	log.Info().Int("count", len(ids)).Msgf("complemented characters")

	return nil
}

func (p *Processor) getCharacterWithFullInfo(ctx context.Context, id int) (*m27r.Character, error) {
	params := &operations.GetCharacterIndividualParams{
		CharacterID: int32(id),
	}
	p.setParams(ctx, params)

	indiv, err := p.mclient.Operations.GetCharacterIndividual(params)
	if err != nil {

		return nil, fmt.Errorf("error fetching character %d: %v", id, err)
	}

	log.Info().Int("id", id).Msg("fetched character with basic info")

	ch := indiv.Payload.Data.Results[0]

	if ch.Comics.Available != ch.Comics.Returned {
		comics, err := p.getCharacterComics(ctx, int32(id), ch.Comics.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching comics for character %d: %v", id, err)
		}

		/*
			skip verification here AS responses differ between
			available returned from /v1/public/characters/{characterId}
			and
			total returned from /v1/public/characters/{characterId}/{comics,events,series,stories}
		*/
		// if ch.Comics.Available != int32(len(comics)) {
		// 	return nil, fmt.Errorf("data missing when fetching comics for character %d: got %d, want %d", id, len(comics), ch.Comics.Available)
		// }

		ch.Comics.Items = comics
		ch.Comics.Returned = ch.Comics.Available
	} else {
		log.Info().Int("id", id).Int32("count", ch.Comics.Available).Msg("character has complete comics")
	}

	if ch.Events.Available != ch.Events.Returned {
		events, err := p.getCharacterEvents(ctx, int32(id), ch.Events.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for character %d: %v", id, err)
		}

		ch.Events.Items = events
		ch.Events.Returned = ch.Events.Available
	} else {
		log.Info().Int("id", id).Int32("count", ch.Events.Available).Msg("character has complete events")
	}

	if ch.Series.Available != ch.Series.Returned {
		series, err := p.getCharacterSeries(ctx, int32(id), ch.Series.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching series for character %d: %v", id, err)
		}

		ch.Series.Items = series
		ch.Series.Returned = ch.Series.Available
	} else {
		log.Info().Int("id", id).Int32("count", ch.Series.Available).Msg("character has complete series")
	}

	if ch.Stories.Available != ch.Stories.Returned {
		stories, err := p.getCharacterStories(ctx, int32(id), ch.Stories.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching stories for character %d: %v", id, err)
		}

		ch.Stories.Items = stories
		ch.Stories.Returned = ch.Stories.Available
	} else {
		log.Info().Int("id", id).Int32("count", ch.Stories.Available).Msg("character has complete stories")
	}

	converted, err := convertCharacter(ch)
	if err != nil {
		return nil, fmt.Errorf("error converting character %d: %v", ch.ID, err)
	}

	return converted, nil
}

func (p *Processor) getCharacterComics(ctx context.Context, id, count int32) ([]*models.ComicSummary, error) {
	var comics []*models.ComicSummary

	comicCh := make(chan *models.ComicSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			params := &operations.GetComicsCharacterCollectionParams{
				CharacterID: id,
				Limit:       &p.limit,
				Offset:      &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetComicsCharacterCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching comics for character %d, offset %d: %v", id, offset, err)
			}

			for _, comic := range col.Payload.Data.Results {
				comicCh <- &models.ComicSummary{Name: comic.Title, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(comic.ID), 10)}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(comicCh)

	for comic := range comicCh {
		comics = append(comics, comic)
	}

	log.Info().Int("count", len(comics)).Int32("character_id", id).Msg("fetched comics for character")

	return comics, nil
}

func (p *Processor) getCharacterEvents(ctx context.Context, id, count int32) ([]*models.EventSummary, error) {
	var events []*models.EventSummary

	eventCh := make(chan *models.EventSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			params := &operations.GetCharacterEventsCollectionParams{
				CharacterID: id,
				Limit:       &p.limit,
				Offset:      &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCharacterEventsCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching events for character %d, offset %d: %v", id, offset, err)
			}

			for _, event := range col.Payload.Data.Results {
				eventCh <- &models.EventSummary{Name: event.Title, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(event.ID), 10)}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(eventCh)

	for event := range eventCh {
		events = append(events, event)
	}

	log.Info().Int("count", len(events)).Int32("character_id", id).Msg("fetched events for character")

	return events, nil
}

func (p *Processor) getCharacterSeries(ctx context.Context, id, count int32) ([]*models.SeriesSummary, error) {
	var series []*models.SeriesSummary

	seriesCh := make(chan *models.SeriesSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			params := &operations.GetCharacterSeriesCollectionParams{
				CharacterID: id,
				Limit:       &p.limit,
				Offset:      &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCharacterSeriesCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching series for character %d, offset %d: %v", id, offset, err)
			}

			for _, series := range col.Payload.Data.Results {
				seriesCh <- &models.SeriesSummary{Name: series.Title, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(series.ID), 10)}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(seriesCh)

	for s := range seriesCh {
		series = append(series, s)
	}

	log.Info().Int("count", len(series)).Int32("character_id", id).Msg("fetched series for character")

	return series, nil
}

func (p *Processor) getCharacterStories(ctx context.Context, id, count int32) ([]*models.StorySummary, error) {
	var stories []*models.StorySummary

	storyCh := make(chan *models.StorySummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			params := &operations.GetCharacterStoryCollectionParams{
				CharacterID: id,
				Limit:       &p.limit,
				Offset:      &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCharacterStoryCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching stories for character %d, offset %d: %v", id, offset, err)
			}

			for _, story := range col.Payload.Data.Results {
				storyCh <- &models.StorySummary{Name: story.Title, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(story.ID), 10)}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(storyCh)

	for story := range storyCh {
		stories = append(stories, story)
	}

	log.Info().Int("count", len(stories)).Int32("character_id", id).Msg("fetched stories for character")

	return stories, nil
}

func convertCharacter(in *models.Character) (*m27r.Character, error) {
	out := &m27r.Character{
		Description: in.Description,
		ID:          in.ID,
		Modified:    in.Modified,
		Name:        in.Name,
		Thumbnail:   strings.Replace(in.Thumbnail.Path+"."+in.Thumbnail.Extension, "http://", "https://", 1),
	}

	for _, item := range in.Comics.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Comics = append(out.Comics, id)
	}

	for _, item := range in.Events.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Events = append(out.Events, id)
	}

	for _, item := range in.Series.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Series = append(out.Series, id)
	}

	for _, item := range in.Stories.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Stories = append(out.Stories, id)
	}

	for _, url := range in.Urls {
		out.URLs = append(out.URLs, &m27r.URL{
			Type: url.Type,
			URL:  strings.Replace(strings.Split(url.URL, "?")[0], "http://", "https://", 1),
		})
	}

	if in.Comics.Available == in.Comics.Returned && in.Events.Available == in.Events.Returned && in.Series.Available == in.Series.Returned && in.Stories.Available == in.Stories.Returned {
		out.Intact = true
	}

	return out, nil
}
