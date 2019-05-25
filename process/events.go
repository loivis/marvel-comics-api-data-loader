package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
	"golang.org/x/sync/errgroup"

	"github.com/loivis/marvel-comics-api-data-loader/marvel/mclient/operations"
	"github.com/loivis/marvel-comics-api-data-loader/marvel/models"
	"github.com/rs/zerolog/log"
)

func (p *Processor) loadEvents(ctx context.Context) error {
	var err error

	err = p.loadAllEventsWithBasicInfo(ctx)
	if err != nil {
		return fmt.Errorf("error loading all events info: %v", err)
	}

	log.Info().Msg("all events loaded")

	// err = p.complementAllEvents(ctx)
	// if err != nil {
	// 	return fmt.Errorf("error complementing all events: %v", err)
	// }

	// log.Info().Msg("all events complemented")

	return nil
}

func (p *Processor) loadAllEventsWithBasicInfo(ctx context.Context) error {
	remote, err := p.getEventCount(ctx)
	if err != nil {
		return fmt.Errorf("error fetching event count: %v", err)
	}
	log.Info().Str("type", "event").Int32("count", remote).Msg("event count from api")

	existing, err := p.store.GetCount("events")
	if err != nil {
		return err
	}
	log.Info().Str("type", "event").Int64("count", existing).Msg("existing event count")

	if int64(remote) == existing {
		log.Info().Int64("local", existing).Int32("remote", remote).Msg("no missing events")
		return nil
	}

	log.Info().Int64("local", existing).Int32("remote", remote).Msg("missing events, reload")

	return p.loadMissingEvents(ctx, int32(existing), remote)
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

			if err := p.store.SaveEvents(events); err != nil {
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
		log.Info().Msg("done")
	}

	return nil
}

func (p *Processor) getEventCount(ctx context.Context) (int32, error) {
	var limit int32 = 1
	params := &operations.GetEventsCollectionParams{
		Limit: &limit,
	}
	p.setParams(ctx, params)

	col, err := p.mclient.Operations.GetEventsCollection(params)
	if err != nil {
		return 0, err
	}

	return col.Payload.Data.Total, nil
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
