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

func (p *Processor) loadSeries(ctx context.Context) error {
	var err error

	/*
		there seems to be duplications in results while paging over /v1/public/series.
		10926 => 10904
		skip load all series based on comparison between api response and local storage.
	*/
	err = p.loadAllSeriesWithBasicInfo(ctx)
	if err != nil {
		return fmt.Errorf("error loading all series info: %v", err)
	}

	log.Info().Msg("all series loaded")

	err = p.complementAllSeries(ctx)
	if err != nil {
		return fmt.Errorf("error complementing all series: %v", err)
	}

	log.Info().Msg("all series complemented")

	return nil
}

func (p *Processor) loadAllSeriesWithBasicInfo(ctx context.Context) error {
	remote, err := p.mclient.GetCount(ctx, maco.TypeSeries)
	if err != nil {
		return fmt.Errorf("error fetching series count: %v", err)
	}
	log.Info().Str("type", "series").Int("count", remote).Msg("series count from api")

	existing, err := p.store.GetCount(ctx, "series")
	if err != nil {
		return err
	}
	log.Info().Str("type", "series").Int("count", existing).Msg("existing series count")

	if remote == existing {
		log.Info().Int("local", existing).Int("remote", remote).Msg("no missing series")
		return nil
	}

	log.Info().Int("local", existing).Int("remote", remote).Msg("missing series, reload")

	return p.loadMissingSeries(ctx, existing, remote)
}

func (p *Processor) loadMissingSeries(ctx context.Context, starting, count int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	seriesCh := make(chan *maco.Series, p.concurrency*p.limit)
	conCh := make(chan struct{}, p.concurrency)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	go func() {
		var series []*maco.Series
		defer func() {
			doneCh <- struct{}{}
		}()

		batchSave := func(series []*maco.Series) error {
			err := retry.Do(func() error {
				return p.store.SaveSeries(ctx, series)
			})

			if err != nil {
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
				series = []*maco.Series{}
			}
		}

		batchSave(series)
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
				log.Info().Int("offset", offset).Msgf("cancelled fetching paged series: %v", err)
				return fmt.Errorf("cancelled fetching paged series limit %d offset %d: %v", p.limit, offset, err)
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int("offset", offset).Msg("ctx already cancelled")
				return nil
			default: // default to avoid blocking
			}

			err := retry.Do(
				func() error {
					series, err := p.mclient.GetSeries(ctx, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
					if err != nil {
						return err
					}

					for _, s := range series {
						converted, err := convertSeries(s)
						if err != nil {
							return err
						}

						seriesCh <- converted
					}

					log.Info().Int("offset", offset).Int("count", len(series)).Msg("fetched paged series")

					return nil
				},
				retry.OnRetry(retryLog(offset)),
				retry.RetryIf(retryIf(offset)),
			)

			if err != nil {
				cancel()
				log.Info().Int("offset", offset).Msgf("cancelled fetching paged series after retry: %v", err)
				return fmt.Errorf("error fetching with limit %d offset %d: (%[3]T) %[3]v", p.limit, offset, err)
			}

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
	ids, err := p.store.IncompleteIDs(ctx, "series")
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

			log.Info().Int("id", id).Msgf("fetched series with full info converted")

			err = p.store.SaveOne(ctx, series)
			if err != nil {
				return fmt.Errorf("error saving series %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("saved series")

			log.Info().Int("id", id).Msgf("complemented series")

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error complementing series: %v", err)
	}

	log.Info().Int("count", len(ids)).Msgf("complemented series")

	return nil
}

func (p *Processor) getSeriesWithFullInfo(ctx context.Context, id int) (*maco.Series, error) {
	series, err := p.mclient.GetSeriesSingle(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching series %d: %v", id, err)
	}

	log.Info().Int("id", id).Msg("fetched series with basic info")

	if series.Characters.Available != series.Characters.Returned {
		chars, err := p.getSeriesCharacters(ctx, id, series.Characters.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching characters for series %d: %v", id, err)
		}

		series.Characters.Items = chars
		series.Characters.Returned = series.Characters.Available
	} else {
		log.Info().Int("id", id).Int("count", series.Characters.Available).Msg("series has complete characters")
	}

	if series.Comics.Available != series.Comics.Returned {
		comics, err := p.getSeriesComics(ctx, id, series.Comics.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching comics for series %d: %v", id, err)
		}

		series.Comics.Items = comics
		series.Comics.Returned = series.Comics.Available
	} else {
		log.Info().Int("id", id).Int("count", series.Comics.Available).Msg("series has complete comics")
	}

	if series.Creators.Available != series.Creators.Returned {
		creators, err := p.getSeriesCreators(ctx, id, series.Creators.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching creators for series %d: %v", id, err)
		}

		series.Creators.Items = creators
		series.Creators.Returned = series.Creators.Available
	} else {
		log.Info().Int("id", id).Int("count", series.Creators.Available).Msg("series has complete creators")
	}

	if series.Events.Available != series.Events.Returned {
		events, err := p.getSeriesEvents(ctx, id, series.Events.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for series %d: %v", id, err)
		}

		series.Events.Items = events
		series.Events.Returned = series.Events.Available
	} else {
		log.Info().Int("id", id).Int("count", series.Events.Available).Msg("series has complete events")
	}

	if series.Stories.Available != series.Stories.Returned {
		stories, err := p.getSeriesStories(ctx, id, series.Stories.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching stories for series %d: %v", id, err)
		}

		series.Stories.Items = stories
		series.Stories.Returned = series.Stories.Available
	} else {
		log.Info().Int("id", id).Int("count", series.Stories.Available).Msg("series has complete stories")
	}

	converted, err := convertSeries(series)
	if err != nil {
		return nil, fmt.Errorf("error converting series %d: %v", series.ID, err)
	}

	return converted, nil
}

func (p *Processor) getSeriesCharacters(ctx context.Context, id, count int) ([]*marvel.CharacterSummary, error) {
	var chars []*marvel.CharacterSummary

	charCh := make(chan *marvel.CharacterSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * i

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			chars, err := p.mclient.GetSeriesCharacters(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching character for series %d, offset %d: %v", id, offset, err)
			}

			for _, char := range chars {
				charCh <- &marvel.CharacterSummary{Name: char.Name, ResourceURI: strconv.Itoa(char.ID)}
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

	log.Info().Int("count", len(chars)).Int("series_id", id).Msg("fetched characters for series")

	return chars, nil
}

func (p *Processor) getSeriesComics(ctx context.Context, id, count int) ([]*marvel.ComicSummary, error) {
	var comics []*marvel.ComicSummary

	comicCh := make(chan *marvel.ComicSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * i

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			comics, err := p.mclient.GetSeriesComics(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching comics for series %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(comics)).Int("series_id", id).Msg("fetched comics for series")

	return comics, nil
}

func (p *Processor) getSeriesCreators(ctx context.Context, id, count int) ([]*marvel.CreatorSummary, error) {
	var creators []*marvel.CreatorSummary

	creatorCh := make(chan *marvel.CreatorSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * i

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			creators, err := p.mclient.GetSeriesCreators(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching creators for series %d, offset %d: %v", id, offset, err)
			}

			for _, creator := range creators {
				creatorCh <- &marvel.CreatorSummary{Name: creator.FullName, ResourceURI: strconv.Itoa(creator.ID)}
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

	log.Info().Int("count", len(creators)).Int("series_id", id).Msg("fetched creators for series")

	return creators, nil
}

func (p *Processor) getSeriesEvents(ctx context.Context, id, count int) ([]*marvel.EventSummary, error) {
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

			events, err := p.mclient.GetSeriesEvents(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching events for series %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(events)).Int("series_id", id).Msg("fetched events for series")

	return events, nil
}

func (p *Processor) getSeriesStories(ctx context.Context, id, count int) ([]*marvel.StorySummary, error) {
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

			stories, err := p.mclient.GetSeriesStories(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching stories for series %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(stories)).Int("series_id", id).Msg("fetched stories for series")

	return stories, nil
}

func convertSeries(in *marvel.Series) (*maco.Series, error) {
	out := &maco.Series{
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

	for _, url := range in.URLs {
		out.URLs = append(out.URLs, &maco.URL{
			Type: url.Type,
			URL:  strings.Replace(strings.Split(url.URL, "?")[0], "http://", "https://", 1),
		})
	}

	if in.Characters.Available == in.Characters.Returned && in.Comics.Available == in.Comics.Returned && in.Creators.Available == in.Creators.Returned && in.Events.Available == in.Events.Returned && in.Stories.Available == in.Stories.Returned {
		out.Intact = true
	}

	return out, nil
}
