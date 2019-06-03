package process

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/avast/retry-go"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/loivis/marvel-comics-api-data-loader/client/marvel"
	"github.com/loivis/marvel-comics-api-data-loader/maco"
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
	remote, err := p.mclient.GetCount(ctx, maco.TypeCharacters)
	if err != nil {
		return fmt.Errorf("error fetching character count: %v", err)
	}
	log.Info().Str("type", "character").Int("count", remote).Msg("character count from api")

	existing, err := p.store.GetCount(ctx, maco.TypeCharacters)
	if err != nil {
		return err
	}
	log.Info().Str("type", "character").Int("count", existing).Msg("existing character count")

	if remote == existing {
		log.Info().Int("local", existing).Int("remote", remote).Msg("no missing characters")
		return nil
	}

	log.Info().Int("local", existing).Int("remote", remote).Msg("missing characters, reload")

	return p.loadMissingCharacters(ctx, existing, remote)
}

func (p *Processor) loadMissingCharacters(ctx context.Context, starting, count int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	charCh := make(chan *maco.Character, p.concurrency*p.limit)
	conCh := make(chan struct{}, p.concurrency)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	go func() {
		var characters []*maco.Character
		defer func() {
			doneCh <- struct{}{}
		}()

		batchSave := func(characters []*maco.Character) error {
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
				characters = []*maco.Character{}
			}
		}

		batchSave(characters)
	}()

	var g errgroup.Group

	for i := starting / p.limit; i < count/p.limit+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * i

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			select {
			case err := <-errCh: // check if any error saving data
				cancel()
				log.Info().Int("offset", offset).Msgf("cancelled fetching paged characters: %v", err)
				return fmt.Errorf("cancelled fetching paged characters limit %d offset %d: %v", p.limit, offset, err)
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int("offset", offset).Msg("ctx already cancelled")
				return nil
			default: // default to avoid blocking
			}

			err := retry.Do(
				func() error {
					chars, err := p.mclient.GetCharacters(ctx, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
					if err != nil {
						return err
					}

					for _, char := range chars {
						converted, err := convertCharacter(char)
						if err != nil {
							return err
						}

						charCh <- converted
					}

					log.Info().Int("offset", offset).Int("count", len(chars)).Msg("fetched paged characters")

					return nil
				},
				retry.OnRetry(retryLog(offset)),
				retry.RetryIf(retryIf(offset)),
			)

			if err != nil {
				cancel()
				log.Info().Int("offset", offset).Msgf("cancelled fetching paged characters after retry: %v", err)
				return fmt.Errorf("error fetching with limit %d offset %d: (%[3]T) %[3]v", p.limit, offset, err)
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
	ids, err := p.store.IncompleteIDs(ctx, maco.TypeCharacters)
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

func (p *Processor) getCharacterWithFullInfo(ctx context.Context, id int) (*maco.Character, error) {
	char, err := p.mclient.GetCharacter(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching character %d: %v", id, err)
	}

	log.Info().Int("id", id).Msg("fetched character with basic info")

	if char.Comics.Available != char.Comics.Returned {
		comics, err := p.getCharacterComics(ctx, id, char.Comics.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching comics for character %d: %v", id, err)
		}

		/*
			skip verification here AS responses differ between
			available returned from /v1/public/characters/{characterId}
			and
			total returned from /v1/public/characters/{characterId}/{comics,events,series,stories}
		*/
		// if char.Comics.Available != int(len(comics)) {
		// 	return nil, fmt.Errorf("data missing when fetching comics for character %d: got %d, want %d", id, len(comics), char.Comics.Available)
		// }

		char.Comics.Items = comics
		char.Comics.Returned = char.Comics.Available
	} else {
		log.Info().Int("id", id).Int("count", char.Comics.Available).Msg("character has complete comics")
	}

	if char.Events.Available != char.Events.Returned {
		events, err := p.getCharacterEvents(ctx, id, char.Events.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for character %d: %v", id, err)
		}

		char.Events.Items = events
		char.Events.Returned = char.Events.Available
	} else {
		log.Info().Int("id", id).Int("count", char.Events.Available).Msg("character has complete events")
	}

	if char.Series.Available != char.Series.Returned {
		series, err := p.getCharacterSeries(ctx, id, char.Series.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching series for character %d: %v", id, err)
		}

		char.Series.Items = series
		char.Series.Returned = char.Series.Available
	} else {
		log.Info().Int("id", id).Int("count", char.Series.Available).Msg("character has complete series")
	}

	if char.Stories.Available != char.Stories.Returned {
		stories, err := p.getCharacterStories(ctx, id, char.Stories.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching stories for character %d: %v", id, err)
		}

		char.Stories.Items = stories
		char.Stories.Returned = char.Stories.Available
	} else {
		log.Info().Int("id", id).Int("count", char.Stories.Available).Msg("character has complete stories")
	}

	converted, err := convertCharacter(char)
	if err != nil {
		return nil, fmt.Errorf("error converting character %d: %v", char.ID, err)
	}

	return converted, nil
}

func (p *Processor) getCharacterComics(ctx context.Context, id, count int) ([]*marvel.ComicSummary, error) {
	var comics []*marvel.ComicSummary

	comicCh := make(chan *marvel.ComicSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < count/p.limit+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * i

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			comics, err := p.mclient.GetCharacterComics(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching comics for character %d, offset %d: %v", id, offset, err)
			}

			for _, comic := range comics {
				comicCh <- &marvel.ComicSummary{Name: comic.Title, ResourceURI: strconv.Itoa(comic.ID)}
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

	log.Info().Int("count", len(comics)).Int("character_id", id).Msg("fetched comics for character")

	return comics, nil
}

func (p *Processor) getCharacterEvents(ctx context.Context, id, count int) ([]*marvel.EventSummary, error) {
	var events []*marvel.EventSummary

	eventCh := make(chan *marvel.EventSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * i

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			events, err := p.mclient.GetCharacterEvents(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching events for character %d, offset %d: %v", id, offset, err)
			}

			for _, event := range events {
				eventCh <- &marvel.EventSummary{Name: event.Title, ResourceURI: strconv.Itoa(event.ID)}
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

	log.Info().Int("count", len(events)).Int("character_id", id).Msg("fetched events for character")

	return events, nil
}

func (p *Processor) getCharacterSeries(ctx context.Context, id, count int) ([]*marvel.SeriesSummary, error) {
	var series []*marvel.SeriesSummary

	seriesCh := make(chan *marvel.SeriesSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * i

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			series, err := p.mclient.GetCharacterSeries(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching series for character %d, offset %d: %v", id, offset, err)
			}

			for _, s := range series {
				seriesCh <- &marvel.SeriesSummary{Name: s.Title, ResourceURI: strconv.Itoa(s.ID)}
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

	log.Info().Int("count", len(series)).Int("character_id", id).Msg("fetched series for character")

	return series, nil
}

func (p *Processor) getCharacterStories(ctx context.Context, id, count int) ([]*marvel.StorySummary, error) {
	var stories []*marvel.StorySummary

	storyCh := make(chan *marvel.StorySummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * i

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			stories, err := p.mclient.GetCharacterStories(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching stories for character %d, offset %d: %v", id, offset, err)
			}

			for _, story := range stories {
				storyCh <- &marvel.StorySummary{Name: story.Title, ResourceURI: strconv.Itoa(story.ID)}
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

	log.Info().Int("count", len(stories)).Int("character_id", id).Msg("fetched stories for character")

	return stories, nil
}

func convertCharacter(in *marvel.Character) (*maco.Character, error) {
	out := &maco.Character{
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

	for _, url := range in.URLs {
		out.URLs = append(out.URLs, &maco.URL{
			Type: url.Type,
			URL:  strings.Replace(strings.Split(url.URL, "?")[0], "http://", "https://", 1),
		})
	}

	if in.Comics.Available == in.Comics.Returned && in.Events.Available == in.Events.Returned && in.Series.Available == in.Series.Returned && in.Stories.Available == in.Stories.Returned {
		out.Intact = true
	}

	return out, nil
}
