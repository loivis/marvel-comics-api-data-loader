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
	"github.com/loivis/marvel-comics-api-data-loader/m27r"
)

func (p *Processor) loadStories(ctx context.Context) error {
	var err error

	err = p.loadAllStoriesWithBasicInfo(ctx)
	if err != nil {
		return fmt.Errorf("error loading all stories info: %v", err)
	}

	log.Info().Msg("all stories loaded")

	err = p.complementAllStories(ctx)
	if err != nil {
		return fmt.Errorf("error complementing all stories: %v", err)
	}

	log.Info().Msg("all stories complemented")

	return nil
}

func (p *Processor) loadAllStoriesWithBasicInfo(ctx context.Context) error {
	remote, err := p.mclient.GetCount(ctx, m27r.TypeStories)
	if err != nil {
		return fmt.Errorf("error fetching story count: %v", err)
	}
	log.Info().Str("type", "story").Int("count", remote).Msg("story count from api")

	existing, err := p.store.GetCount(ctx, "stories")
	if err != nil {
		return err
	}
	log.Info().Str("type", "story").Int("count", existing).Msg("existing story count")

	if remote == existing {
		log.Info().Int("local", existing).Int("remote", remote).Msg("no missing stories")
		return nil
	}

	log.Info().Int("local", existing).Int("remote", remote).Msg("missing stories, reload")

	return p.loadMissingStories(ctx, existing, remote)
}

func (p *Processor) loadMissingStories(ctx context.Context, starting, count int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	storyCh := make(chan *m27r.Story, p.concurrency*p.limit)
	conCh := make(chan struct{}, p.concurrency)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	go func() {
		var stories []*m27r.Story
		defer func() {
			doneCh <- struct{}{}
		}()

		batchSave := func(stories []*m27r.Story) error {
			err := retry.Do(func() error {
				return p.store.SaveStories(ctx, stories)
			})

			if err != nil {
				return err
			}

			log.Info().Int("count", len(stories)).Msg("batch saved stories")

			return nil
		}

		for story := range storyCh {
			stories = append(stories, story)

			if len(stories) >= p.storeBatch {
				if err := batchSave(stories); err != nil {
					errCh <- err
					break
				}
				stories = []*m27r.Story{}
			}
		}

		batchSave(stories)
	}()

	var g errgroup.Group

	for i := int(starting / p.limit); i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * i

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			select {
			case err := <-errCh: // check if any error saving data
				cancel()
				log.Info().Int("offset", offset).Msgf("cancelled fetching paged stories: %v", err)
				return fmt.Errorf("cancelled fetching paged stories limit %d offset %d: %v", p.limit, offset, err)
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int("offset", offset).Msg("ctx already cancelled")
				return nil
			default: // default to avoid blocking
			}

			err := retry.Do(
				func() error {
					stories, err := p.mclient.GetStories(ctx, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
					if err != nil {
						return err
					}

					for _, story := range stories {
						converted, err := convertStory(story)
						if err != nil {
							return err
						}

						storyCh <- converted
					}

					log.Info().Int("offset", offset).Int("count", len(stories)).Msg("fetched paged stories")

					return nil
				},
				retry.OnRetry(retryLog(offset)),
				retry.RetryIf(retryIf(offset)),
			)

			if err != nil {
				cancel()
				log.Info().Int("offset", offset).Msg("cancelled fetching paged stories")
				return fmt.Errorf("error fetching with limit %d offset %d: (%[3]T) %[3]v", p.limit, offset, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	close(storyCh)

	select {
	case <-doneCh:
		log.Info().Msg("fetched all missing stories with basic info")
	}

	return nil
}

func (p *Processor) complementAllStories(ctx context.Context) error {
	ids, err := p.store.IncompleteIDs(ctx, "stories")
	if err != nil {
		return fmt.Errorf("error get imcomplete story ids: %v", err)
	}

	if len(ids) == 0 {
		log.Info().Msg("no incomplete story")
		return nil
	}

	log.Info().Int("count", len(ids)).Msg("fetched incomplete story ids")

	var g errgroup.Group

	conCh := make(chan struct{}, p.concurrency)

	for _, id := range ids {
		conCh <- struct{}{}

		id := id

		g.Go(func() error {
			defer func() {
				<-conCh
			}()

			story, err := p.getStoryWithFullInfo(ctx, id)
			if err != nil {
				return fmt.Errorf("error fetching story %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("fetched story with full info converted")

			err = p.store.SaveOne(ctx, story)
			if err != nil {
				return fmt.Errorf("error saving story %d: %v", id, err)
			}

			log.Info().Int("id", id).Msgf("saved story")

			log.Info().Int("id", id).Msgf("complemented story")

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error complementing stories: %v", err)
	}

	log.Info().Int("count", len(ids)).Msgf("complemented stories")

	return nil
}

func (p *Processor) getStoryWithFullInfo(ctx context.Context, id int) (*m27r.Story, error) {
	story, err := p.mclient.GetStory(ctx, id)
	if err != nil {

		return nil, fmt.Errorf("error fetching story %d: %v", id, err)
	}

	log.Info().Int("id", id).Msg("fetched story with basic info")

	if story.Characters.Available != story.Characters.Returned {
		chars, err := p.getStoryCharacters(ctx, id, story.Characters.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching characters for story %d: %v", id, err)
		}

		story.Characters.Items = chars
		story.Characters.Returned = story.Characters.Available
	} else {
		log.Info().Int("id", id).Int("count", story.Characters.Available).Msg("story has complete characters")
	}

	if story.Comics.Available != story.Comics.Returned {
		comics, err := p.getStoryComics(ctx, id, story.Comics.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching comics for story %d: %v", id, err)
		}

		story.Comics.Items = comics
		story.Comics.Returned = story.Comics.Available
	} else {
		log.Info().Int("id", id).Int("count", story.Comics.Available).Msg("story has complete comics")
	}

	if story.Creators.Available != story.Creators.Returned {
		creators, err := p.getStoryCreators(ctx, id, story.Creators.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching creators for story %d: %v", id, err)
		}

		story.Creators.Items = creators
		story.Creators.Returned = story.Creators.Available
	} else {
		log.Info().Int("id", id).Int("count", story.Creators.Available).Msg("story has complete creators")
	}

	if story.Events.Available != story.Events.Returned {
		events, err := p.getStoryEvents(ctx, id, story.Events.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for story %d: %v", id, err)
		}

		story.Events.Items = events
		story.Events.Returned = story.Events.Available
	} else {
		log.Info().Int("id", id).Int("count", story.Events.Available).Msg("story has complete events")
	}

	if story.Series.Available != story.Series.Returned {
		series, err := p.getStorySeries(ctx, id, story.Series.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching series for story %d: %v", id, err)
		}

		story.Series.Items = series
		story.Series.Returned = story.Series.Available
	} else {
		log.Info().Int("id", id).Int("count", story.Series.Available).Msg("story has complete series")
	}

	converted, err := convertStory(story)
	if err != nil {
		return nil, fmt.Errorf("error converting story %d: %v", story.ID, err)
	}

	return converted, nil
}

func (p *Processor) getStoryCharacters(ctx context.Context, id, count int) ([]*marvel.CharacterSummary, error) {
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

			chars, err := p.mclient.GetStoryCharacters(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching character for story %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(chars)).Int("story_id", id).Msg("fetched characters for story")

	return chars, nil
}

func (p *Processor) getStoryComics(ctx context.Context, id, count int) ([]*marvel.ComicSummary, error) {
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

			comics, err := p.mclient.GetStoryComics(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching comics for story %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(comics)).Int("story_id", id).Msg("fetched comics for story")

	return comics, nil
}

func (p *Processor) getStoryCreators(ctx context.Context, id, count int) ([]*marvel.CreatorSummary, error) {
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

			creators, err := p.mclient.GetStoryCreators(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching creators for story %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(creators)).Int("story_id", id).Msg("fetched creators for story")

	return creators, nil
}

func (p *Processor) getStoryEvents(ctx context.Context, id, count int) ([]*marvel.EventSummary, error) {
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

			events, err := p.mclient.GetStoryEvents(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching events for story %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(events)).Int("story_id", id).Msg("fetched events for story")

	return events, nil
}

func (p *Processor) getStorySeries(ctx context.Context, id, count int) ([]*marvel.SeriesSummary, error) {
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

			series, err := p.mclient.GetStorySeries(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching series for story %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(series)).Int("story_id", id).Msg("fetched series for story")

	return series, nil
}

func convertStory(in *marvel.Story) (*m27r.Story, error) {
	out := &m27r.Story{
		Description: in.Description,
		ID:          in.ID,
		Modified:    in.Modified,
		Title:       in.Title,
		Type:        in.Type,
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

	if in.Originalissue != nil {
		id, err := idFromURL(in.Originalissue.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", in.Originalissue.ResourceURI, err)
		}
		out.Originalissue = id
	}

	for _, item := range in.Series.Items {
		id, err := idFromURL(item.ResourceURI)
		if err != nil {
			return nil, fmt.Errorf("error get id from %q: %v", item.ResourceURI, err)
		}

		out.Series = append(out.Series, id)
	}

	if in.Thumbnail != nil {
		out.Thumbnail = strings.Replace(in.Thumbnail.Path+"."+in.Thumbnail.Extension, "http://", "https://", 1)
	}

	if in.Characters.Available == in.Characters.Returned && in.Comics.Available == in.Comics.Returned && in.Creators.Available == in.Creators.Returned && in.Events.Available == in.Events.Returned && in.Series.Available == in.Series.Returned {
		out.Intact = true
	}

	return out, nil
}
