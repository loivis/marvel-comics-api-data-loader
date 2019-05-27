package process

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
	"github.com/loivis/marvel-comics-api-data-loader/marvel/mclient/operations"
	"github.com/loivis/marvel-comics-api-data-loader/marvel/models"
)

func (p *Processor) loadSeries(ctx context.Context) error {
	var err error

	/*
		there seems to be duplications in results while paging over /v1/public/series.
		10926 => 10904
		skip load all comics based on comparison between api response and local storage.
	*/
	// err = p.loadAllSeriesWithBasicInfo(ctx)
	// if err != nil {
	// 	return fmt.Errorf("error loading all series info: %v", err)
	// }

	// log.Info().Msg("all series loaded")

	err = p.complementAllSeries(ctx)
	if err != nil {
		return fmt.Errorf("error complementing all series: %v", err)
	}

	log.Info().Msg("all series complemented")

	return nil
}

func (p *Processor) loadAllSeriesWithBasicInfo(ctx context.Context) error {
	remote, err := p.getSeriesCount(ctx)
	if err != nil {
		return fmt.Errorf("error fetching series count: %v", err)
	}
	log.Info().Str("type", "series").Int32("count", remote).Msg("series count from api")

	existing, err := p.store.GetCount("series")
	if err != nil {
		return err
	}
	log.Info().Str("type", "series").Int64("count", existing).Msg("existing series count")

	if int64(remote) == existing {
		log.Info().Int64("local", existing).Int32("remote", remote).Msg("no missing series")
		return nil
	}

	log.Info().Int64("local", existing).Int32("remote", remote).Msg("missing series, reload")

	return p.loadMissingSeries(ctx, int32(existing), remote)
}

func (p *Processor) getSeriesCount(ctx context.Context) (int32, error) {
	var limit int32 = 1
	params := &operations.GetSeriesCollectionParams{
		Limit: &limit,
	}
	p.setParams(ctx, params)

	col, err := p.mclient.Operations.GetSeriesCollection(params)
	if err != nil {
		return 0, err
	}

	return col.Payload.Data.Total, nil
}

func (p *Processor) loadMissingSeries(ctx context.Context, starting, count int32) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	seriesCh := make(chan *m27r.Series, int32(p.concurrency)*p.limit)
	conCh := make(chan struct{}, p.concurrency)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	go func() {
		var series []*m27r.Series
		defer func() {
			doneCh <- struct{}{}
		}()

		batchSave := func(series []*m27r.Series) error {
			if err := p.store.SaveSeries(series); err != nil {
				return err
			}

			log.Info().Int("count", len(series)).Msg("batch saved series")

			return nil
		}

		for s := range seriesCh {
			series = append(series, s)

			if len(series) >= p.storeBatch {
				if err := batchSave(series); err != nil {
					errCh <- err
					break
				}
				series = []*m27r.Series{}
			}
		}

		batchSave(series)
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
				log.Info().Int32("offset", offset).Msgf("cancelled fetching paged series: %v", err)
				return fmt.Errorf("cancelled fetching paged series limit %d offset %d: %v", p.limit, offset, err)
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("offset", offset).Msg("ctx already cancelled")
				return nil
			default: // default to avoid blocking
			}

			params := &operations.GetSeriesCollectionParams{
				Limit:  &p.limit,
				Offset: &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetSeriesCollection(params)
			if err != nil {
				cancel()
				log.Info().Int32("offset", offset).Msg("cancelled fetching paged series")
				return fmt.Errorf("error fetching with limit %d offset %d: %v", p.limit, offset, err)
			}

			for _, res := range col.Payload.Data.Results {
				series, err := convertSeries(res)
				if err != nil {
					return err
				}

				seriesCh <- series
			}

			log.Info().Int32("offset", offset).Int32("count", col.Payload.Data.Count).Msg("fetched paged series")

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	close(seriesCh)

	select {
	case <-doneCh:
		log.Info().Msg("fetched all missing series with basic info")
	}

	return nil
}

func (p *Processor) complementAllSeries(ctx context.Context) error {
	ids, err := p.store.IncompleteIDs("series")
	if err != nil {
		return fmt.Errorf("error get imcomplete series ids: %v", err)
	}

	if len(ids) == 0 {
		log.Info().Msg("no incomplete series")
		return nil
	}

	log.Info().Int("count", len(ids)).Msg("fetched incomplete series ids")

	var g errgroup.Group

	conCh := make(chan struct{}, p.concurrency)

	for _, id := range ids {
		conCh <- struct{}{}

		id := id

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			series, err := p.getSeriesWithFullInfo(ctx, id)
			if err != nil {
				return fmt.Errorf("error fetching series %d: %v", id, err)
			}

			log.Info().Int32("id", id).Msgf("fetched series with full info converted")

			err = p.store.SaveOne(series)
			if err != nil {
				return fmt.Errorf("error saving series %d: %v", id, err)
			}

			log.Info().Int32("id", id).Msgf("saved series")

			log.Info().Int32("id", id).Msgf("complemented series")

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error complementing series: %v", err)
	}

	log.Info().Int("count", len(ids)).Msgf("complemented series")

	return nil
}

func (p *Processor) getSeriesWithFullInfo(ctx context.Context, id int32) (*m27r.Series, error) {
	params := &operations.GetSeriesIndividualParams{
		SeriesID: id,
	}
	p.setParams(ctx, params)

	indiv, err := p.mclient.Operations.GetSeriesIndividual(params)
	if err != nil {
		return nil, fmt.Errorf("error fetching series %d: %v", id, err)
	}

	log.Info().Int32("id", id).Msg("fetched series with basic info")

	series := indiv.Payload.Data.Results[0]

	if series.Characters.Available != series.Characters.Returned {
		chars, err := p.getSeriesCharacters(ctx, id, series.Characters.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching characters for series %d: %v", id, err)
		}

		series.Characters.Items = chars
		series.Characters.Returned = series.Characters.Available
	} else {
		log.Info().Int32("id", id).Int32("count", series.Characters.Available).Msg("series has complete characters")
	}

	if series.Comics.Available != series.Comics.Returned {
		comics, err := p.getSeriesComics(ctx, id, series.Comics.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching comics for series %d: %v", id, err)
		}

		series.Comics.Items = comics
		series.Comics.Returned = series.Comics.Available
	} else {
		log.Info().Int32("id", id).Int32("count", series.Comics.Available).Msg("series has complete comics")
	}

	if series.Creators.Available != series.Creators.Returned {
		creators, err := p.getSeriesCreators(ctx, id, series.Creators.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching creators for series %d: %v", id, err)
		}

		series.Creators.Items = creators
		series.Creators.Returned = series.Creators.Available
	} else {
		log.Info().Int32("id", id).Int32("count", series.Creators.Available).Msg("series has complete creators")
	}

	if series.Events.Available != series.Events.Returned {
		events, err := p.getSeriesEvents(ctx, id, series.Events.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for series %d: %v", id, err)
		}

		series.Events.Items = events
		series.Events.Returned = series.Events.Available
	} else {
		log.Info().Int32("id", id).Int32("count", series.Events.Available).Msg("series has complete events")
	}

	if series.Stories.Available != series.Stories.Returned {
		stories, err := p.getSeriesStories(ctx, id, series.Stories.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching stories for series %d: %v", id, err)
		}

		series.Stories.Items = stories
		series.Stories.Returned = series.Stories.Available
	} else {
		log.Info().Int32("id", id).Int32("count", series.Stories.Available).Msg("series has complete stories")
	}

	e, err := convertSeries(series)
	if err != nil {
		return nil, fmt.Errorf("error converting series %d: %v", series.ID, err)
	}

	return e, nil
}

func (p *Processor) getSeriesCharacters(ctx context.Context, id, count int32) ([]*models.CharacterSummary, error) {
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

			params := &operations.GetSeriesCharacterWrapperParams{
				SeriesID: id,
				Limit:    &p.limit,
				Offset:   &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetSeriesCharacterWrapper(params)
			if err != nil {
				return fmt.Errorf("error fetching character for series %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(chars)).Int32("series_id", id).Msg("fetched characters for series")

	return chars, nil
}

func (p *Processor) getSeriesComics(ctx context.Context, id, count int32) ([]*models.ComicSummary, error) {
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

			params := &operations.GetComicsCollectionBySeriesIDParams{
				SeriesID: id,
				Limit:    &p.limit,
				Offset:   &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetComicsCollectionBySeriesID(params)
			if err != nil {
				return fmt.Errorf("error fetching comics for series %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(comics)).Int32("series_id", id).Msg("fetched comics for series")

	return comics, nil
}

func (p *Processor) getSeriesCreators(ctx context.Context, id, count int32) ([]*models.CreatorSummary, error) {
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

			params := &operations.GetCreatorCollectionBySeriesIDParams{
				SeriesID: id,
				Limit:    &p.limit,
				Offset:   &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCreatorCollectionBySeriesID(params)
			if err != nil {
				return fmt.Errorf("error fetching creators for series %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(creators)).Int32("series_id", id).Msg("fetched creators for series")

	return creators, nil
}

func (p *Processor) getSeriesEvents(ctx context.Context, id, count int32) ([]*models.EventSummary, error) {
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

			params := &operations.GetEventsCollectionBySeriesIDParams{
				SeriesID: id,
				Limit:    &p.limit,
				Offset:   &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetEventsCollectionBySeriesID(params)
			if err != nil {
				return fmt.Errorf("error fetching events for series %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(events)).Int32("series_id", id).Msg("fetched events for series")

	return events, nil
}

func (p *Processor) getSeriesStories(ctx context.Context, id, count int32) ([]*models.StorySummary, error) {
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

			params := &operations.GetSeriesStoryCollectionParams{
				SeriesID: id,
				Limit:    &p.limit,
				Offset:   &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetSeriesStoryCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching stories for series %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(stories)).Int32("series_id", id).Msg("fetched stories for series")

	return stories, nil
}

func convertSeries(in *models.Series) (*m27r.Series, error) {
	out := &m27r.Series{
		Description: in.Description,
		EndYear:     in.EndYear,
		ID:          in.ID,
		Modified:    in.Modified,
		Rating:      in.Rating,
		StartYear:   in.StartYear,
		Thumbnail:   strings.Replace(in.Thumbnail.Path+"."+in.Thumbnail.Extension, "http://", "https://", 1),
		Title:       in.Title,
	}

	for _, item := range in.Characters.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Characters = append(out.Characters, id)
	}

	for _, item := range in.Comics.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Comics = append(out.Comics, id)
	}

	for _, item := range in.Creators.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Creators = append(out.Creators, id)
	}

	for _, item := range in.Events.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Events = append(out.Events, id)
	}

	if in.Next != nil {
		id, err := idFromURL(in.Next.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", in.Next.ResourceURI, err)
		}
		out.Next = id
	}

	if in.Previous != nil {
		id, err := idFromURL(in.Previous.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", in.Previous.ResourceURI, err)
		}
		out.Next = id
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

	if in.Characters.Available == in.Characters.Returned && in.Comics.Available == in.Comics.Returned && in.Creators.Available == in.Creators.Returned && in.Events.Available == in.Events.Returned && in.Stories.Available == in.Stories.Returned {
		out.Intact = true
	}

	return out, nil
}
