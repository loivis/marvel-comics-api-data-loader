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

func (p *Processor) loadCreators(ctx context.Context) error {
	var err error

	/*
		there are creators with only id and wrong-type fields while paging over /v1/public/creators.
		6213 => ...
		skip load all creators based on comparison between api response and local storage.
	*/
	err = p.loadAllCreatorsWithBasicInfo(ctx)
	if err != nil {
		return fmt.Errorf("error loading all creators info: %v", err)
	}

	log.Info().Msg("all creators loaded")

	err = p.complementAllCreators(ctx)
	if err != nil {
		return fmt.Errorf("error complementing all creators: %v", err)
	}

	log.Info().Msg("all creators complemented")

	return nil
}

func (p *Processor) loadAllCreatorsWithBasicInfo(ctx context.Context) error {
	remote, err := p.mclient.GetCount(ctx, maco.TypeCreators)
	if err != nil {
		return fmt.Errorf("error fetching creator count: %v", err)
	}
	log.Info().Str("type", "creator").Int("count", remote).Msg("creator count from api")

	existing, err := p.store.GetCount(ctx, "creators")
	if err != nil {
		return err
	}
	log.Info().Str("type", "creator").Int("count", existing).Msg("existing creator count")

	if remote == existing {
		log.Info().Int("local", existing).Int("remote", remote).Msg("no missing creators")
		return nil
	}

	log.Info().Int("local", existing).Int("remote", remote).Msg("missing creators, reload")

	return p.loadMissingCreators(ctx, existing, remote)
}

func (p *Processor) loadMissingCreators(ctx context.Context, starting, count int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	creatorCh := make(chan *maco.Creator, p.concurrency*p.limit)
	conCh := make(chan struct{}, p.concurrency)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	go func() {
		var creators []*maco.Creator
		defer func() {
			doneCh <- struct{}{}
		}()

		batchSave := func(creators []*maco.Creator) error {
			err := retry.Do(func() error {
				return p.store.SaveCreators(ctx, creators)
			})

			if err != nil {
				return err
			}

			log.Info().Int("count", len(creators)).Msg("batch saved creators")

			return nil
		}

		for creator := range creatorCh {
			creators = append(creators, creator)

			if len(creators) >= p.storeBatch {
				if err := batchSave(creators); err != nil {
					errCh <- err
					break
				}
				creators = []*maco.Creator{}
			}
		}

		batchSave(creators)
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
				log.Info().Int("offset", offset).Msgf("cancelled fetching paged creators: %v", err)
				return fmt.Errorf("cancelled fetching paged creators limit %d offset %d: %v", p.limit, offset, err)
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int("offset", offset).Msg("ctx already cancelled")
				return nil
			default: // default to avoid blocking
			}

			err := retry.Do(
				func() error {
					creators, err := p.mclient.GetCreators(ctx, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
					if err != nil {
						return err
					}

					for _, creator := range creators {
						converted, err := convertCreator(creator)
						if err != nil {
							return err
						}

						creatorCh <- converted
					}

					log.Info().Int("offset", offset).Int("count", len(creators)).Msg("fetched paged creators")

					return nil
				},
				retry.OnRetry(retryLog(offset)),
				retry.RetryIf(retryIf(offset)),
			)

			if err != nil {
				cancel()
				log.Info().Int("offset", offset).Msgf("cancelled fetching paged creators after retry: %v", err)
				return fmt.Errorf("error fetching with limit %d offset %d: (%[3]T) %[3]v", p.limit, offset, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	close(creatorCh)

	select {
	case <-doneCh:
		log.Info().Msg("fetched all missing creators with basic info")
	}

	return nil
}

func (p *Processor) complementAllCreators(ctx context.Context) error {
	ids, err := p.store.IncompleteIDs(ctx, "creators")
	if err != nil {
		return fmt.Errorf("error get imcomplete creator ids: %v", err)
	}

	if len(ids) == 0 {
		log.Info().Msg("no incomplete creator")
		return nil
	}

	log.Info().Int("count", len(ids)).Msg("fetched incomplete creator ids")

	var g errgroup.Group

	conCh := make(chan struct{}, p.concurrency)

	for _, id := range ids {
		conCh <- struct{}{}

		id := id

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			creator, err := p.getCreatorWithFullInfo(ctx, id)
			if err != nil {
				return fmt.Errorf("error fetching creator %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("fetched creator with full info converted")

			err = p.store.SaveOne(ctx, creator)
			if err != nil {
				return fmt.Errorf("error saving creator %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("saved creator")

			log.Info().Int("id", id).Msgf("complemented creator")

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error complementing creators: %v", err)
	}

	log.Info().Int("count", len(ids)).Msgf("complemented creators")

	return nil
}

func (p *Processor) getCreatorWithFullInfo(ctx context.Context, id int) (*maco.Creator, error) {
	creator, err := p.mclient.GetCreator(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching creator %d: %v", id, err)
	}

	log.Info().Int("id", id).Msg("fetched creator with basic info")

	/*
		skip verification all below AS responses may differ between
		available returned from /v1/public/creators/{creatorId}
		and
		total returned from /v1/public/creators/{creatorId}/{comics/events/series/stories}
	*/

	if creator.Comics.Available != creator.Comics.Returned {
		comics, err := p.getCreatorComics(ctx, id, creator.Comics.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching comics for creator %d: %v", id, err)
		}

		creator.Comics.Items = comics
		creator.Comics.Returned = creator.Comics.Available
	} else {
		log.Info().Int("id", id).Int("count", creator.Comics.Available).Msg("creator has complete comics")
	}

	if creator.Events.Available != creator.Events.Returned {
		events, err := p.getCreatorEvents(ctx, id, creator.Events.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for creator %d: %v", id, err)
		}

		creator.Events.Items = events
		creator.Events.Returned = creator.Events.Available
	} else {
		log.Info().Int("id", id).Int("count", creator.Events.Available).Msg("creator has complete events")
	}

	if creator.Series.Available != creator.Series.Returned {
		series, err := p.getCreatorSeries(ctx, id, creator.Series.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching series for creator %d: %v", id, err)
		}

		creator.Series.Items = series
		creator.Series.Returned = creator.Series.Available
	} else {
		log.Info().Int("id", id).Int("count", creator.Series.Available).Msg("creator has complete series")
	}

	if creator.Stories.Available != creator.Stories.Returned {
		stories, err := p.getCreatorStories(ctx, id, creator.Stories.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching stories for creator %d: %v", id, err)
		}

		creator.Stories.Items = stories
		creator.Stories.Returned = creator.Stories.Available
	} else {
		log.Info().Int("id", id).Int("count", creator.Stories.Available).Msg("creator has complete stories")
	}

	converted, err := convertCreator(creator)
	if err != nil {
		return nil, fmt.Errorf("error converting creator %d: %v", creator.ID, err)
	}

	return converted, nil
}

func (p *Processor) getCreatorComics(ctx context.Context, id, count int) ([]*marvel.ComicSummary, error) {
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

			comics, err := p.mclient.GetCreatorComics(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching comics for creator %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(comics)).Int("creator_id", id).Msg("fetched comics for creator")

	return comics, nil
}

func (p *Processor) getCreatorEvents(ctx context.Context, id, count int) ([]*marvel.EventSummary, error) {
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

			events, err := p.mclient.GetCreatorEvents(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching events for creator %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(events)).Int("creator_id", id).Msg("fetched events for creator")

	return events, nil
}

func (p *Processor) getCreatorSeries(ctx context.Context, id, count int) ([]*marvel.SeriesSummary, error) {
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

			series, err := p.mclient.GetCreatorSeries(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching series for creator %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(series)).Int("creator_id", id).Msg("fetched series for creator")

	return series, nil
}

func (p *Processor) getCreatorStories(ctx context.Context, id, count int) ([]*marvel.StorySummary, error) {
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

			stories, err := p.mclient.GetCreatorStories(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching stories for creator %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(stories)).Int("creator_id", id).Msg("fetched stories for creator")

	return stories, nil
}

func convertCreator(in *marvel.Creator) (*maco.Creator, error) {
	out := &maco.Creator{
		FirstName:  in.FirstName,
		FullName:   in.FullName,
		ID:         in.ID,
		LastName:   in.LastName,
		MiddleName: in.MiddleName,
		Modified:   in.Modified,
		Suffix:     in.Suffix,
		Thumbnail:  strings.Replace(in.Thumbnail.Path+"."+in.Thumbnail.Extension, "http://", "https://", 1),
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
