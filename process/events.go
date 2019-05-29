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

func (p *Processor) loadEvents(ctx context.Context) error {
	var err error

	err = p.loadAllEventsWithBasicInfo(ctx)
	if err != nil {
		return fmt.Errorf("error loading all events info: %v", err)
	}

	log.Info().Msg("all events loaded")

	err = p.complementAllEvents(ctx)
	if err != nil {
		return fmt.Errorf("error complementing all events: %v", err)
	}

	log.Info().Msg("all events complemented")

	return nil
}

func (p *Processor) loadAllEventsWithBasicInfo(ctx context.Context) error {
	remote, err := p.getEventCount(ctx)
	if err != nil {
		return fmt.Errorf("error fetching event count: %v", err)
	}
	log.Info().Str("type", "event").Int("count", remote).Msg("event count from api")

	existing, err := p.store.GetCount(ctx, "events")
	if err != nil {
		return err
	}
	log.Info().Str("type", "event").Int("count", existing).Msg("existing event count")

	if int(remote) == existing {
		log.Info().Int("local", existing).Int("remote", remote).Msg("no missing events")
		return nil
	}

	log.Info().Int("local", existing).Int("remote", remote).Msg("missing events, reload")

	return p.loadMissingEvents(ctx, int32(existing), int32(remote))
}

func (p *Processor) getEventCount(ctx context.Context) (int, error) {
	var limit int32 = 1
	params := &operations.GetEventsCollectionParams{
		Limit: &limit,
	}
	p.setParams(ctx, params)

	col, err := p.mclient.Operations.GetEventsCollection(params)
	if err != nil {
		return 0, err
	}

	return int(col.Payload.Data.Total), nil
}

func (p *Processor) loadMissingEvents(ctx context.Context, starting, count int32) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eventCh := make(chan *m27r.Event, int32(p.concurrency)*p.limit)
	conCh := make(chan struct{}, p.concurrency)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	go func() {
		var events []*m27r.Event
		defer func() {
			doneCh <- struct{}{}
		}()

		batchSave := func(events []*m27r.Event) error {
			if err := p.store.SaveEvents(ctx, events); err != nil {
				return err
			}

			log.Info().Int("count", len(events)).Msg("batch saved events")

			return nil
		}

		for event := range eventCh {
			events = append(events, event)

			if len(events) >= p.storeBatch {
				if err := batchSave(events); err != nil {
					errCh <- err
					break
				}
				events = []*m27r.Event{}
			}
		}

		batchSave(events)
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
				log.Info().Int32("offset", offset).Msgf("cancelled fetching paged events: %v", err)
				return fmt.Errorf("cancelled fetching paged events limit %d offset %d: %v", p.limit, offset, err)
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("offset", offset).Msg("ctx already cancelled")
				return nil
			default: // default to avoid blocking
			}

			params := &operations.GetEventsCollectionParams{
				Limit:  &p.limit,
				Offset: &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetEventsCollection(params)
			if err != nil {
				cancel()
				log.Info().Int32("offset", offset).Msg("cancelled fetching paged events")
				return fmt.Errorf("error fetching with limit %d offset %d: %v", p.limit, offset, err)
			}

			for _, res := range col.Payload.Data.Results {
				event, err := convertEvent(res)
				if err != nil {
					return err
				}

				eventCh <- event
			}

			log.Info().Int32("offset", offset).Int32("count", col.Payload.Data.Count).Msg("fetched paged events")

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	close(eventCh)

	select {
	case <-doneCh:
		log.Info().Msg("fetched all missing events with basic info")
	}

	return nil
}

func (p *Processor) complementAllEvents(ctx context.Context) error {
	ids, err := p.store.IncompleteIDs(ctx, "events")
	if err != nil {
		return fmt.Errorf("error get imcomplete event ids: %v", err)
	}

	if len(ids) == 0 {
		log.Info().Msg("no incomplete event")
		return nil
	}

	log.Info().Int("count", len(ids)).Msg("fetched incomplete event ids")

	var g errgroup.Group

	conCh := make(chan struct{}, p.concurrency)

	for _, id := range ids {
		conCh <- struct{}{}

		id := id

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			event, err := p.getEventWithFullInfo(ctx, id)
			if err != nil {
				return fmt.Errorf("error fetching event %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("fetched event with full info converted")

			err = p.store.SaveOne(ctx, event)
			if err != nil {
				return fmt.Errorf("error saving event %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("saved event")

			log.Info().Int("id", id).Msgf("complemented event")

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error complementing events: %v", err)
	}

	log.Info().Int("count", len(ids)).Msgf("complemented events")

	return nil
}

func (p *Processor) getEventWithFullInfo(ctx context.Context, id int) (*m27r.Event, error) {
	params := &operations.GetEventIndividualParams{
		EventID: int32(id),
	}
	p.setParams(ctx, params)

	indiv, err := p.mclient.Operations.GetEventIndividual(params)
	if err != nil {

		return nil, fmt.Errorf("error fetching event %d: %v", id, err)
	}

	log.Info().Int("id", id).Msg("fetched event with basic info")

	event := indiv.Payload.Data.Results[0]

	if event.Characters.Available != event.Characters.Returned {
		chars, err := p.getEventCharacters(ctx, int32(id), event.Characters.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching characters for event %d: %v", id, err)
		}

		event.Characters.Items = chars
		event.Characters.Returned = event.Characters.Available
	} else {
		log.Info().Int("id", id).Int32("count", event.Characters.Available).Msg("event has complete characters")
	}

	if event.Comics.Available != event.Comics.Returned {
		comics, err := p.getEventComics(ctx, int32(id), event.Comics.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching comics for event %d: %v", id, err)
		}

		event.Comics.Items = comics
		event.Comics.Returned = event.Comics.Available
	} else {
		log.Info().Int("id", id).Int32("count", event.Comics.Available).Msg("event has complete comics")
	}

	if event.Creators.Available != event.Creators.Returned {
		creators, err := p.getEventCreators(ctx, int32(id), event.Creators.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching creators for event %d: %v", id, err)
		}

		event.Creators.Items = creators
		event.Creators.Returned = event.Creators.Available
	} else {
		log.Info().Int("id", id).Int32("count", event.Creators.Available).Msg("event has complete creators")
	}

	if event.Series.Available != event.Series.Returned {
		series, err := p.getEventSeries(ctx, int32(id), event.Series.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching series for event %d: %v", id, err)
		}

		event.Series.Items = series
		event.Series.Returned = event.Series.Available
	} else {
		log.Info().Int("id", id).Int32("count", event.Series.Available).Msg("event has complete series")
	}

	if event.Stories.Available != event.Stories.Returned {
		stories, err := p.getEventStories(ctx, int32(id), event.Stories.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching stories for event %d: %v", id, err)
		}

		event.Stories.Items = stories
		event.Stories.Returned = event.Stories.Available
	} else {
		log.Info().Int("id", id).Int32("count", event.Stories.Available).Msg("event has complete stories")
	}

	converted, err := convertEvent(event)
	if err != nil {
		return nil, fmt.Errorf("error converting event %d: %v", event.ID, err)
	}

	return converted, nil
}

func (p *Processor) getEventCharacters(ctx context.Context, id, count int32) ([]*models.CharacterSummary, error) {
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

			params := &operations.GetEventCharacterCollectionParams{
				EventID: id,
				Limit:   &p.limit,
				Offset:  &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetEventCharacterCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching character for event %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(chars)).Int32("event_id", id).Msg("fetched characters for event")

	return chars, nil
}

func (p *Processor) getEventComics(ctx context.Context, id, count int32) ([]*models.ComicSummary, error) {
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

			params := &operations.GetComicsCollectionByEventIDParams{
				EventID: id,
				Limit:   &p.limit,
				Offset:  &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetComicsCollectionByEventID(params)
			if err != nil {
				return fmt.Errorf("error fetching comics for event %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(comics)).Int32("event_id", id).Msg("fetched comics for event")

	return comics, nil
}

func (p *Processor) getEventCreators(ctx context.Context, id, count int32) ([]*models.CreatorSummary, error) {
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

			params := &operations.GetCreatorCollectionByEventIDParams{
				EventID: id,
				Limit:   &p.limit,
				Offset:  &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCreatorCollectionByEventID(params)
			if err != nil {
				return fmt.Errorf("error fetching creators for event %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(creators)).Int32("event_id", id).Msg("fetched creators for event")

	return creators, nil
}

func (p *Processor) getEventSeries(ctx context.Context, id, count int32) ([]*models.SeriesSummary, error) {
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

			params := &operations.GetEventSeriesCollectionParams{
				EventID: id,
				Limit:   &p.limit,
				Offset:  &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetEventSeriesCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching series for event %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(series)).Int32("event_id", id).Msg("fetched series for event")

	return series, nil
}

func (p *Processor) getEventStories(ctx context.Context, id, count int32) ([]*models.StorySummary, error) {
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

			params := &operations.GetEventStoryCollectionParams{
				EventID: id,
				Limit:   &p.limit,
				Offset:  &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetEventStoryCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching stories for event %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(stories)).Int32("event_id", id).Msg("fetched stories for event")

	return stories, nil
}

func convertEvent(in *models.Event) (*m27r.Event, error) {
	out := &m27r.Event{
		Description: in.Description,
		End:         in.End,
		ID:          in.ID,
		Modified:    in.Modified,
		Start:       in.Start,
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

	if in.Characters.Available == in.Characters.Returned && in.Comics.Available == in.Comics.Returned && in.Creators.Available == in.Creators.Returned && in.Series.Available == in.Series.Returned && in.Stories.Available == in.Stories.Returned {
		out.Intact = true
	}

	return out, nil
}
