package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
	"github.com/loivis/marvel-comics-api-data-loader/marvel/mclient/operations"
	"github.com/loivis/marvel-comics-api-data-loader/marvel/models"
)

func (p *Processor) loadComics(ctx context.Context) error {
	err := p.loadAllComicsWithBasicInfo(ctx)
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
	/*
		should skip count comparison between api response and local storage.
		there seems to be duplications in results while paging over /v1/public/comics.
		MAYBE NOT? It might because of change on limit/offset.
	*/
	remote, err := p.getComicCount(ctx)
	if err != nil {
		return fmt.Errorf("error fetching comic count: %v", err)
	}
	log.Info().Str("type", "comic").Int32("count", remote).Msg("fetched")

	existing, err := p.store.GetCount("comics")
	if err != nil {
		return err
	}
	log.Info().Str("type", "comic").Int64("count", existing).Msg("existing comics")

	if int64(remote) == existing {
		log.Info().Int64("local", existing).Int32("remote", remote).Msg("no missing comics")
		return nil
	}

	log.Info().Int64("local", existing).Int32("remote", remote).Msg("missing comics, reload")
	comics, err := p.getAllComics(ctx, int32(existing), remote)
	if err != nil {
		return fmt.Errorf("error getting all comics: %v", err)
	}

	log.Info().Int("count", len(comics)).Msg("fetched all comics")

	err = p.store.SaveComics(comics)
	if err != nil {
		return fmt.Errorf("error storing comics: %v", err)
	}

	return nil
}

func (p *Processor) getComicCount(ctx context.Context) (int32, error) {
	var limit int32 = 1
	params := &operations.GetComicsCollectionParams{
		Limit: &limit,
	}
	p.setParams(ctx, params)

	col, err := p.mclient.Operations.GetComicsCollection(params)
	if err != nil {
		return 0, err
	}

	return col.Payload.Data.Total, nil
}

// assume that results in response from /v1/public/comics is ordered list
// so we can start from where we don't have data.
func (p *Processor) getAllComics(ctx context.Context, starting int32, count int32) ([]*m27r.Comic, error) {
	var comics []*m27r.Comic

	comicCh := make(chan *m27r.Comic, count)
	conCh := make(chan struct{}, p.concurrency)

	var g errgroup.Group

	for i := int(starting / p.limit); i < int(count/p.limit)+1; i++ {
		conCh <- struct{}{}
		offset := p.limit * int32(i)
		g.Go(func() error {
			params := &operations.GetComicsCollectionParams{
				Limit:  &p.limit,
				Offset: &offset,
			}
			p.setParams(ctx, params)

			col, err := p.mclient.Operations.GetComicsCollection(params)
			if err != nil {
				return fmt.Errorf("error fetching with limit %d offset %d: %v", p.limit, offset, err)
			}

			for _, res := range col.Payload.Data.Results {
				comic, err := convertComic(res)
				if err != nil {
					return err
				}

				comicCh <- comic
			}

			log.Info().Int32("offset", offset).Int32("count", col.Payload.Data.Count).Msg("fetched paged comics")

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

	return comics, nil
}

func (p *Processor) complementAllComics(ctx context.Context) error {
	return nil
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
