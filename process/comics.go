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
	remote, err := p.mclient.GetCount(ctx, maco.TypeComics)
	if err != nil {
		return fmt.Errorf("error fetching comic count: %v", err)
	}
	log.Info().Str("type", "comic").Int("count", remote).Msg("comic count from api")

	existing, err := p.store.GetCount(ctx, "comics")
	if err != nil {
		return err
	}
	log.Info().Str("type", "comic").Int("count", existing).Msg("existing comic count")

	if remote == existing {
		log.Info().Int("local", existing).Int("remote", remote).Msg("no missing comics")
		return nil
	}

	log.Info().Int("local", existing).Int("remote", remote).Msg("missing comics, reload")

	return p.loadMissingComics(ctx, existing, remote)
}

func (p *Processor) loadMissingComics(ctx context.Context, starting, count int) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	comicCh := make(chan *maco.Comic, p.concurrency*p.limit)
	conCh := make(chan struct{}, p.concurrency)
	errCh := make(chan error, 1)
	doneCh := make(chan struct{})

	go func() {
		var comics []*maco.Comic
		defer func() {
			doneCh <- struct{}{}
		}()

		batchSave := func(comics []*maco.Comic) error {
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
				comics = []*maco.Comic{}
			}
		}

		batchSave(comics)
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
				log.Info().Int("offset", offset).Msgf("cancelled fetching paged comics: %v", err)
				return fmt.Errorf("cancelled fetching paged comics limit %d offset %d: %v", p.limit, offset, err)
			case <-ctx.Done(): // Check if ctx was cancelled in other goroutine
				log.Info().Int("offset", offset).Msg("ctx already cancelled")
				return nil
			default: // default to avoid blocking
			}

			err := retry.Do(
				func() error {
					comics, err := p.mclient.GetComics(ctx, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
					if err != nil {
						return err
					}

					for _, comic := range comics {
						converted, err := convertComic(comic)
						if err != nil {
							return err
						}

						comicCh <- converted
					}

					log.Info().Int("offset", offset).Int("count", len(comics)).Msg("fetched paged comics")

					return nil
				},
				retry.OnRetry(retryLog(offset)),
				retry.RetryIf(retryIf(offset)),
			)

			if err != nil {
				cancel()
				log.Info().Int("offset", offset).Msgf("cancelled fetching paged comics after retry: %v", err)
				return fmt.Errorf("error fetching with limit %d offset %d: (%[3]T) %[3]v", p.limit, offset, err)
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

func (p *Processor) getComicWithFullInfo(ctx context.Context, id int) (*maco.Comic, error) {
	comic, err := p.mclient.GetComic(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching comic %d: %v", id, err)
	}

	log.Info().Int("id", id).Msg("fetched comic with basic info")

	if comic.Characters.Available != comic.Characters.Returned {
		chars, err := p.getComicCharacters(ctx, id, comic.Characters.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching characters for comic %d: %v", id, err)
		}

		comic.Characters.Items = chars
		comic.Characters.Returned = comic.Characters.Available
	} else {
		log.Info().Int("id", id).Int("count", comic.Characters.Available).Msg("comic has complete characters")
	}

	if comic.Creators.Available != comic.Creators.Returned {
		creators, err := p.getComicCreators(ctx, id, comic.Creators.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching creators for comic %d: %v", id, err)
		}

		comic.Creators.Items = creators
		comic.Creators.Returned = comic.Creators.Available
	} else {
		log.Info().Int("id", id).Int("count", comic.Creators.Available).Msg("comic has complete creators")
	}

	if comic.Events.Available != comic.Events.Returned {
		events, err := p.getComicEvents(ctx, id, comic.Events.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for comic %d: %v", id, err)
		}

		comic.Events.Items = events
		comic.Events.Returned = comic.Events.Available
	} else {
		log.Info().Int("id", id).Int("count", comic.Events.Available).Msg("comic has complete events")
	}

	if comic.Stories.Available != comic.Stories.Returned {
		stories, err := p.getComicStories(ctx, id, comic.Stories.Available)
		if err != nil {
			return nil, fmt.Errorf("error fetching stories for comic %d: %v", id, err)
		}

		comic.Stories.Items = stories
		comic.Stories.Returned = comic.Stories.Available
	} else {
		log.Info().Int("id", id).Int("count", comic.Stories.Available).Msg("comic has complete stories")
	}

	converted, err := convertComic(comic)
	if err != nil {
		return nil, fmt.Errorf("error converting comic %d: %v", comic.ID, err)
	}

	return converted, nil
}

func (p *Processor) getComicCharacters(ctx context.Context, id, count int) ([]*marvel.CharacterSummary, error) {
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

			chars, err := p.mclient.GetComicCharacters(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching character for comic %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(chars)).Int("comic_id", id).Msg("fetched characters for comic")

	return chars, nil
}

func (p *Processor) getComicCreators(ctx context.Context, id, count int) ([]*marvel.CreatorSummary, error) {
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

			creators, err := p.mclient.GetComicCreators(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching creators for comic %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(creators)).Int("comic_id", id).Msg("fetched creators for comic")

	return creators, nil
}

func (p *Processor) getComicEvents(ctx context.Context, id, count int) ([]*marvel.EventSummary, error) {
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

			events, err := p.mclient.GetComicEvents(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching events for comic %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(events)).Int("comic_id", id).Msg("fetched events for comic")

	return events, nil
}

func (p *Processor) getComicStories(ctx context.Context, id, count int) ([]*marvel.StorySummary, error) {
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

			stories, err := p.mclient.GetComicStories(ctx, id, &marvel.Params{Limit: p.limit, Offset: offset, OrderBy: "modified"})
			if err != nil {
				return fmt.Errorf("error fetching stories for comic %d, offset %d: %v", id, offset, err)
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

	log.Info().Int("count", len(stories)).Int("comic_id", id).Msg("fetched stories for comic")

	return stories, nil
}

func convertComic(in *marvel.Comic) (*maco.Comic, error) {
	out := &maco.Comic{
		Description:        in.Description,
		DiamondCode:        in.DiamondCode,
		DigitalID:          in.DigitalID,
		EAN:                in.Ean,
		Format:             in.Format,
		ID:                 in.ID,
		ISBN:               in.ISBN,
		ISSN:               in.ISSN,
		IssueNumber:        in.IssueNumber,
		Modified:           in.Modified,
		PageCount:          in.PageCount,
		Thumbnail:          strings.Replace(in.Thumbnail.Path+"."+in.Thumbnail.Extension, "http://", "https://", 1),
		Title:              in.Title,
		UPC:                in.UPC,
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
		out.Dates = append(out.Dates, &maco.ComicDate{
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
		out.Prices = append(out.Prices, &maco.ComicPrice{
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
		out.TextObjects = append(out.TextObjects, &maco.TextObject{
			Language: item.Language,
			Text:     item.Text,
			Type:     item.Type,
		})
	}

	for _, url := range in.URLs {
		out.URLs = append(out.URLs, &maco.URL{
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
