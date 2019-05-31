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

func (p *Processor) loadComics(ctx context.Context) error {
	var err error

	/*
	   there seems to be duplications in results while paging over /v1/public/comics.
	   44228 => 41868
	   skip load all comics based on comparison between api response and local storage.
	*/
	err = p.loadAllComicsWithBasicInfo(ctx)
	if err != nil {
		return fmt.Errorf("error loading all comics info: %v", err)
	}

	log.Info().Msg("all comics loaded")

	err = p.complementAllComics(ctx)
	if err != nil {
		return fmt.Errorf("error complementing all comics: %v", err)
	}

	log.Info().Msg("all comics complemented")

	return nil
}

func (p *Processor) loadAllComicsWithBasicInfo(ctx context.Context) error {
	remote, err := p.getComicCount(ctx)
	if err != nil {
		return fmt.Errorf("error fetching comic count: %v", err)
	}
	log.Info().Str("type", "comic").Int("count", remote).Msg("comic count from api")

	existing, err := p.store.GetCount(ctx, "comics")
	if err != nil {
		return err
	}
	log.Info().Str("type", "comic").Int("count", existing).Msg("existing comic count")

	if int(remote) == existing {
		log.Info().Int("local", existing).Int("remote", remote).Msg("no missing comics")
		return nil
	}

	log.Info().Int("local", existing).Int("remote", remote).Msg("missing comics, reload")

	return p.loadMissingComics(ctx, int32(existing), int32(remote))
}

func (p *Processor) getComicCount(ctx context.Context) (int, error) {
	var limit int32 = 1
	params := &operations.GetComicsCollectionParams{
		Limit: &limit,
	}
	p.setParams(ctx, params)

	col, err := p.mclient.Operations.GetComicsCollection(params)
	if err != nil {
		return 0, err
	}

	return int(col.Payload.Data.Total), nil
}

func (p *Processor) loadMissingComics(ctx context.Context, starting, count int32) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	comicCh := make(chan *m27r.Comic, int32(p.concurrency)*p.limit)
	conCh := make(chan struct{}, p.concurrency)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	go func() {
		var comics []*m27r.Comic
		defer func() {
			doneCh <- struct{}{}
		}()

		batchSave := func(comics []*m27r.Comic) error {
			err := retry.Do(func() error {
				return p.store.SaveComics(ctx, comics)
			})

			if err != nil {
				return err
			}

			log.Info().Int("count", len(comics)).Msg("batch saved comics")

			return nil
		}

		for comic := range comicCh {
			comics = append(comics, comic)

			if len(comics) >= p.storeBatch {
				if err := batchSave(comics); err != nil {
					errCh <- err
					break
				}
				comics = []*m27r.Comic{}
			}
		}

		batchSave(comics)
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
				log.Info().Int32("offset", offset).Msgf("cancelled fetching paged comics: %v", err)
				return fmt.Errorf("cancelled fetching paged comics limit %d offset %d: %v", p.limit, offset, err)
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("offset", offset).Msg("ctx already cancelled")
				return nil
			default: // default to avoid blocking
			}

			params := &operations.GetComicsCollectionParams{
				Limit:  &p.limit,
				Offset: &offset,
			}
			p.setParams(ctx, params)

			err := retry.Do(
				func() error {
					col, err := p.mclient.Operations.GetComicsCollection(params)
					if err != nil {
						return err
					}

					for _, res := range col.Payload.Data.Results {
						comic, err := convertComic(res)
						if err != nil {
							return err
						}

						comicCh <- comic
					}

					log.Info().Int32("offset", offset).Int32("count", col.Payload.Data.Count).Msg("fetched paged comics")

					return nil
				},
				retry.OnRetry(retryLog(offset)),
				retry.RetryIf(retryIf(offset)),
			)

			if err != nil {
				cancel()
				log.Info().Int32("offset", offset).Msg("cancelled fetching paged comics")
				return fmt.Errorf("error fetching with limit %d offset %d: (%T) %v", p.limit, offset, err, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	close(comicCh)

	select {
	case <-doneCh:
		log.Info().Msg("fetched all missing comics with basic info")
	}

	return nil
}

func (p *Processor) complementAllComics(ctx context.Context) error {
	ids, err := p.store.IncompleteIDs(ctx, "comics")
	if err != nil {
		return fmt.Errorf("error get imcomplete comic ids: %v", err)
	}

	if len(ids) == 0 {
		log.Info().Msg("no incomplete comic")
		return nil
	}

	log.Info().Int("count", len(ids)).Msg("fetched incomplete comic ids")

	var g errgroup.Group

	conCh := make(chan struct{}, p.concurrency)

	for _, id := range ids {
		conCh <- struct{}{}

		id := id

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			comic, err := p.getComicWithFullInfo(ctx, id)
			if err != nil {
				return fmt.Errorf("error fetching comic %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("fetched comic with full info converted")

			err = p.store.SaveOne(ctx, comic)
			if err != nil {
				return fmt.Errorf("error saving comic %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("saved comic")

			log.Info().Int("id", id).Msgf("complemented comic")

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error complementing comics: %v", err)
	}

	log.Info().Int("count", len(ids)).Msgf("complemented comics")

	return nil
}

func (p *Processor) getComicWithFullInfo(ctx context.Context, id int) (*m27r.Comic, error) {
	params := &operations.GetComicIndividualParams{
		ComicID: int32(id),
	}
	p.setParams(ctx, params)

	indiv, err := p.mclient.Operations.GetComicIndividual(params)
	if err != nil {

		return nil, fmt.Errorf("error fetching comic %d: %v", id, err)
	}

	log.Info().Int("id", id).Msg("fetched comic with basic info")

	comic := indiv.Payload.Data.Results[0]

	if comic.Characters.Available != comic.Characters.Returned {
		chars, err := p.getComicCharacters(ctx, int32(id), comic.Characters.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching characters for comic %d: %v", id, err)
		}

		comic.Characters.Items = chars
		comic.Characters.Returned = comic.Characters.Available
	} else {
		log.Info().Int("id", id).Int32("count", comic.Characters.Available).Msg("comic has complete characters")
	}

	if comic.Creators.Available != comic.Creators.Returned {
		creators, err := p.getComicCreators(ctx, int32(id), comic.Creators.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching creators for comic %d: %v", id, err)
		}

		comic.Creators.Items = creators
		comic.Creators.Returned = comic.Creators.Available
	} else {
		log.Info().Int("id", id).Int32("count", comic.Creators.Available).Msg("comic has complete creators")
	}

	if comic.Events.Available != comic.Events.Returned {
		events, err := p.getComicEvents(ctx, int32(id), comic.Events.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for comic %d: %v", id, err)
		}

		comic.Events.Items = events
		comic.Events.Returned = comic.Events.Available
	} else {
		log.Info().Int("id", id).Int32("count", comic.Events.Available).Msg("comic has complete events")
	}

	if comic.Stories.Available != comic.Stories.Returned {
		stories, err := p.getComicStories(ctx, int32(id), comic.Stories.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching stories for comic %d: %v", id, err)
		}

		comic.Stories.Items = stories
		comic.Stories.Returned = comic.Stories.Available
	} else {
		log.Info().Int("id", id).Int32("count", comic.Stories.Available).Msg("comic has complete stories")
	}

	converted, err := convertComic(comic)
	if err != nil {
		return nil, fmt.Errorf("error converting comic %d: %v", comic.ID, err)
	}

	return converted, nil
}

func (p *Processor) getComicCharacters(ctx context.Context, id, count int32) ([]*models.CharacterSummary, error) {
	var chars []*models.CharacterSummary

	charCh := make(chan *models.CharacterSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			params := &operations.GetComicCharacterCollectionParams{
				ComicID: id,
				Limit:   &p.limit,
				Offset:  &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetComicCharacterCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching character for comic %d, offset %d: %v", id, offset, err)
			}

			for _, char := range col.Payload.Data.Results {
				charCh <- &models.CharacterSummary{Name: char.Name, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(char.ID), 10)}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(charCh)

	for char := range charCh {
		chars = append(chars, char)
	}

	log.Info().Int("count", len(chars)).Int32("comic_id", id).Msg("fetched characters for comic")

	return chars, nil
}

func (p *Processor) getComicCreators(ctx context.Context, id, count int32) ([]*models.CreatorSummary, error) {
	var creators []*models.CreatorSummary

	creatorCh := make(chan *models.CreatorSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			params := &operations.GetCreatorCollectionByComicIDParams{
				ComicID: id,
				Limit:   &p.limit,
				Offset:  &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCreatorCollectionByComicID(params)
			if err != nil {
				return fmt.Errorf("error fetching creators for comic %d, offset %d: %v", id, offset, err)
			}

			for _, creator := range col.Payload.Data.Results {
				creatorCh <- &models.CreatorSummary{Name: creator.FullName, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(creator.ID), 10)}
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	close(creatorCh)

	for creator := range creatorCh {
		creators = append(creators, creator)
	}

	log.Info().Int("count", len(creators)).Int32("comic_id", id).Msg("fetched creators for comic")

	return creators, nil
}

func (p *Processor) getComicEvents(ctx context.Context, id, count int32) ([]*models.EventSummary, error) {
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

			params := &operations.GetIssueEventsCollectionParams{
				ComicID: id,
				Limit:   &p.limit,
				Offset:  &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetIssueEventsCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching events for comic %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(events)).Int32("comic_id", id).Msg("fetched events for comic")

	return events, nil
}

func (p *Processor) getComicStories(ctx context.Context, id, count int32) ([]*models.StorySummary, error) {
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

			params := &operations.GetComicStoryCollectionParams{
				ComicID: id,
				Limit:   &p.limit,
				Offset:  &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetComicStoryCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching stories for comic %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(stories)).Int32("comic_id", id).Msg("fetched stories for comic")

	return stories, nil
}

func convertComic(in *models.Comic) (*m27r.Comic, error) {
	out := &m27r.Comic{
		Description:        in.Description,
		DigitalID:          in.DigitalID,
		EAN:                in.Ean,
		Format:             in.Format,
		ID:                 in.ID,
		ISSN:               in.Issn,
		IssueNumber:        in.IssueNumber,
		Modified:           in.Modified,
		PageCount:          in.PageCount,
		Thumbnail:          strings.Replace(in.Thumbnail.Path+"."+in.Thumbnail.Extension, "http://", "https://", 1),
		Title:              in.Title,
		UPC:                in.Upc,
		VariantDescription: in.VariantDescription,
	}

	id, err := idFromURL(in.Series.ResourceURI)
	if err != nil {
		return nil, fmt.Errorf("error get id from %q: %v", in.Series.ResourceURI, err)
	}
	out.SeriesID = id

	for _, item := range in.Characters.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Characters = append(out.Characters, id)
	}

	for _, item := range in.CollectedIssues {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.CollectedIssues = append(out.CollectedIssues, id)
	}

	for _, item := range in.Collections {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Collections = append(out.Collections, id)
	}

	for _, item := range in.Creators.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Creators = append(out.Creators, id)
	}

	for _, item := range in.Dates {
		out.Dates = append(out.Dates, &m27r.ComicDate{
			Date: item.Date,
			Type: item.Type,
		})
	}

	for _, item := range in.Events.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Events = append(out.Events, id)
	}

	for _, item := range in.Prices {
		out.Prices = append(out.Prices, &m27r.ComicPrice{
			Price: item.Price,
			Type:  item.Type,
		})
	}

	for _, item := range in.Stories.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Stories = append(out.Stories, id)
	}

	for _, item := range in.TextObjects {
		out.TextObjects = append(out.TextObjects, &m27r.TextObject{
			Language: item.Language,
			Text:     item.Text,
			Type:     item.Type,
		})
	}

	for _, url := range in.Urls {
		out.URLs = append(out.URLs, &m27r.URL{
			Type: url.Type,
			URL:  strings.Replace(strings.Split(url.URL, "?")[0], "http://", "https://", 1),
		})
	}

	for _, item := range in.Variants {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Variants = append(out.Variants, id)
	}

	if in.Characters.Available == in.Characters.Returned && in.Creators.Available == in.Creators.Returned && in.Events.Available == in.Events.Returned && in.Stories.Available == in.Stories.Returned {
		out.Intact = true
	}

	return out, nil
}
