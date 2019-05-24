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

func (p *Processor) loadCreators(ctx context.Context) error {
	var err error

	/*
		there are creators with only id and wrong-type fields while paging over /v1/public/creators.
		6213 => ...
		skip load all creators based on comparison between api response and local storage.
	*/
	// err = p.loadAllCreatorsWithBasicInfo(ctx)
	// if err != nil {
	// 	return fmt.Errorf("error loading all creators info: %v", err)
	// }

	// log.Info().Msg("all creators loaded")

	err = p.complementAllCreators(ctx)
	if err != nil {
		return fmt.Errorf("error complementing all creators: %v", err)
	}

	log.Info().Msg("all creators complemented")

	return nil
}

func (p *Processor) loadAllCreatorsWithBasicInfo(ctx context.Context) error {
	remote, err := p.getCreatorCount(ctx)
	if err != nil {
		return fmt.Errorf("error fetching creator count: %v", err)
	}
	log.Info().Str("type", "creator").Int32("count", remote).Msg("fetched")

	existing, err := p.store.GetCount("creators")
	if err != nil {
		return err
	}
	log.Info().Str("type", "creator").Int64("count", existing).Msg("existing creators")

	if int64(remote) == existing {
		log.Info().Int64("local", existing).Int32("remote", remote).Msg("no missing creators")
		return nil
	}

	log.Info().Int64("local", existing).Int32("remote", remote).Msg("missing creators, reload")
	creators, err := p.getAllCreators(ctx, int32(existing), remote)
	if err != nil {
		return fmt.Errorf("error getting all creators: %v", err)
	}

	log.Info().Int("count", len(creators)).Msg("fetched all creators")

	err = p.store.SaveCreators(creators)
	if err != nil {
		return fmt.Errorf("error storing creators: %v", err)
	}

	return nil
}

func (p *Processor) getCreatorCount(ctx context.Context) (int32, error) {
	var limit int32 = 1
	params := &operations.GetCreatorCollectionParams{
		Limit: &limit,
	}
	p.setParams(ctx, params)

	col, err := p.mclient.Operations.GetCreatorCollection(params)
	if err != nil {
		return 0, err
	}

	return col.Payload.Data.Total, nil
}

// assume that results in response from /v1/public/creators is an ordered list
// so we can start from where we don't have data.
func (p *Processor) getAllCreators(ctx context.Context, starting int32, count int32) ([]*m27r.Creator, error) {
	var creators []*m27r.Creator

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	creatorCh := make(chan *m27r.Creator, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := int(starting / p.limit); i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)
		g.Go(func() error {
			switch offset {
			case 200, 5200: // containing wrong-type field creators: 9551, 10669
				<-conCh
				return nil
			}

			select {
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("offset", offset).Msg("ctx already cancelled")
				<-conCh
				return nil
			default: // default to avoid blocking
			}

			params := &operations.GetCreatorCollectionParams{
				Limit:  &p.limit,
				Offset: &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCreatorCollection(params)
			if err != nil {
				cancel()
				log.Info().Int32("offset", offset).Msg("cancelled fetching paged creators")
				<-conCh
				return fmt.Errorf("error fetching with limit %d offset %d: %v", p.limit, offset, err)
			}

			for _, res := range col.Payload.Data.Results {
				creator, err := convertCreator(res)
				if err != nil {
					<-conCh
					return err
				}

				creatorCh <- creator
			}

			log.Info().Int32("offset", offset).Int32("count", col.Payload.Data.Count).Msg("fetched paged creators")

			<-conCh
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

	return creators, nil
}

func (p *Processor) complementAllCreators(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ids, err := p.store.IncompleteIDs("creators")
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
			select {
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("id", id).Msg("ctx already cancelled")
				<-conCh
				return nil
			default: // default to avoid blocking
			}

			creator, err := p.getCreatorWithFullInfo(ctx, id)
			if err != nil {
				cancel()
				log.Info().Int32("id", id).Msg("cancelled getting creator with full info")
				<-conCh
				return fmt.Errorf("error fetching creator %d: %v", id, err)
			}

			log.Info().Int32("id", id).Msgf("fetched creator with full info converted")

			err = p.store.SaveOne(creator)
			if err != nil {
				<-conCh
				return fmt.Errorf("error saving creator %d: %v", id, err)
			}

			log.Info().Int32("id", id).Msgf("saved creator")

			log.Info().Int32("id", id).Msgf("complemented creator")

			<-conCh

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error complementing creators: %v", err)
	}

	log.Info().Int("count", len(ids)).Msgf("complemented creators")

	return nil
}

func (p *Processor) getCreatorWithFullInfo(ctx context.Context, id int32) (*m27r.Creator, error) {
	params := &operations.GetCreatorIndividualParams{
		CreatorID: id,
	}
	p.setParams(ctx, params)

	indiv, err := p.mclient.Operations.GetCreatorIndividual(params)
	if err != nil {

		return nil, fmt.Errorf("error fetching creator %d: %v", id, err)
	}

	log.Info().Int32("id", id).Msg("fetched creator with basic info")

	/*
		skip verification all below AS responses may differ between
		available returned from /v1/public/creators/{creatorId}
		and
		total returned from /v1/public/creators/{creatorId}/{comics/events/series/stories}
	*/

	creator := indiv.Payload.Data.Results[0]

	if creator.Comics.Available != creator.Comics.Returned {
		chars, err := p.getCreatorComics(ctx, id, creator.Comics.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching comics for creator %d: %v", id, err)
		}

		creator.Comics.Items = chars
		creator.Comics.Returned = creator.Comics.Available
	} else {
		log.Info().Int32("id", id).Int32("count", creator.Comics.Available).Msg("creator has complete comics")
	}

	if creator.Events.Available != creator.Events.Returned {
		events, err := p.getCreatorEvents(ctx, id, creator.Events.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for creator %d: %v", id, err)
		}

		creator.Events.Items = events
		creator.Events.Returned = creator.Events.Available
	} else {
		log.Info().Int32("id", id).Int32("count", creator.Events.Available).Msg("creator has complete events")
	}

	if creator.Series.Available != creator.Series.Returned {
		series, err := p.getCreatorSeries(ctx, id, creator.Series.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching series for creator %d: %v", id, err)
		}

		creator.Series.Items = series
		creator.Series.Returned = creator.Series.Available
	} else {
		log.Info().Int32("id", id).Int32("count", creator.Series.Available).Msg("creator has complete series")
	}

	if creator.Stories.Available != creator.Stories.Returned {
		stories, err := p.getCreatorStories(ctx, id, creator.Stories.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching stories for creator %d: %v", id, err)
		}

		creator.Stories.Items = stories
		creator.Stories.Returned = creator.Stories.Available
	} else {
		log.Info().Int32("id", id).Int32("count", creator.Stories.Available).Msg("creator has complete stories")
	}

	c, err := convertCreator(creator)
	if err != nil {
		return nil, fmt.Errorf("error converting creator %d: %v", creator.ID, err)
	}

	return c, nil
}

func (p *Processor) getCreatorComics(ctx context.Context, id, count int32) ([]*models.ComicSummary, error) {
	var comics []*models.ComicSummary

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	comicCh := make(chan *models.ComicSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			select {
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("creator_id", id).Int32("offset", offset).Msg("ctx already cancelled")
				<-conCh
				return nil
			default: // default to avoid blocking
			}

			params := &operations.GetComicsCollectionByCreatorIDParams{
				CreatorID: id,
				Limit:     &p.limit,
				Offset:    &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetComicsCollectionByCreatorID(params)
			if err != nil {
				cancel()
				log.Info().Int32("creator_id", id).Int32("offset", offset).Msg("cancelled getting comics for creator")
				<-conCh
				return fmt.Errorf("error fetching comics for creator %d, offset %d: %v", id, offset, err)
			}

			for _, comic := range col.Payload.Data.Results {
				comicCh <- &models.ComicSummary{Name: comic.Title, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(comic.ID), 10)}
			}

			<-conCh
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

	log.Info().Int("count", len(comics)).Int32("creator_id", id).Msg("fetched comics for creator")

	return comics, nil
}

func (p *Processor) getCreatorEvents(ctx context.Context, id, count int32) ([]*models.EventSummary, error) {
	var events []*models.EventSummary

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	eventCh := make(chan *models.EventSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			select {
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("creator_id", id).Int32("offset", offset).Msg("ctx already cancelled")
				<-conCh
				return nil
			default: // default to avoid blocking
			}

			params := &operations.GetCreatorEventsCollectionParams{
				CreatorID: id,
				Limit:     &p.limit,
				Offset:    &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCreatorEventsCollection(params)
			if err != nil {
				cancel()
				log.Info().Int32("creator_id", id).Int32("offset", offset).Msg("cancelled getting events for creator")
				<-conCh
				return fmt.Errorf("error fetching events for creator %d, offset %d: %v", id, offset, err)
			}

			for _, event := range col.Payload.Data.Results {
				eventCh <- &models.EventSummary{Name: event.Title, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(event.ID), 10)}
			}

			<-conCh
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

	log.Info().Int("count", len(events)).Int32("creator_id", id).Msg("fetched events for creator")

	return events, nil
}

func (p *Processor) getCreatorSeries(ctx context.Context, id, count int32) ([]*models.SeriesSummary, error) {
	var series []*models.SeriesSummary

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	seriesCh := make(chan *models.SeriesSummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			select {
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("creator_id", id).Int32("offset", offset).Msg("ctx already cancelled")
				<-conCh
				return nil
			default: // default to avoid blocking
			}

			params := &operations.GetCreatorSeriesCollectionParams{
				CreatorID: id,
				Limit:     &p.limit,
				Offset:    &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCreatorSeriesCollection(params)
			if err != nil {
				cancel()
				log.Info().Int32("creator_id", id).Int32("offset", offset).Msg("cancelled getting series for creator")
				<-conCh
				return fmt.Errorf("error fetching series for creator %d, offset %d: %v", id, offset, err)
			}

			for _, series := range col.Payload.Data.Results {
				seriesCh <- &models.SeriesSummary{Name: series.Title, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(series.ID), 10)}
			}

			<-conCh
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

	log.Info().Int("count", len(series)).Int32("creator_id", id).Msg("fetched series for creator")

	return series, nil
}

func (p *Processor) getCreatorStories(ctx context.Context, id, count int32) ([]*models.StorySummary, error) {
	var stories []*models.StorySummary

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	storyCh := make(chan *models.StorySummary, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)

		g.Go(func() error {
			switch {
			case id == 2041 && offset == 0: // containing wrong-type field stories: 30688, 44568
				<-conCh
				return nil
			}

			select {
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int32("creator_id", id).Int32("offset", offset).Msg("ctx already cancelled")
				<-conCh
				return nil
			default: // default to avoid blocking
			}

			params := &operations.GetCreatorStoryCollectionParams{
				CreatorID: id,
				Limit:     &p.limit,
				Offset:    &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetCreatorStoryCollection(params)
			if err != nil {
				cancel()
				log.Info().Int32("creator_id", id).Int32("offset", offset).Msg("cancelled getting stories for creator")
				<-conCh
				return fmt.Errorf("error fetching stories for creator %d, offset %d: %v", id, offset, err)
			}

			for _, story := range col.Payload.Data.Results {
				storyCh <- &models.StorySummary{Name: story.Title, ResourceURI: "fake-prefix/" + strconv.FormatInt(int64(story.ID), 10)}
			}

			<-conCh
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

	log.Info().Int("count", len(stories)).Int32("creator_id", id).Msg("fetched stories for creator")

	return stories, nil
}

func convertCreator(in *models.Creator) (*m27r.Creator, error) {
	out := &m27r.Creator{
		FirtName:   in.FirstName,
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
