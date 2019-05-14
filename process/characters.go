package process

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/loivis/mcapi-loader/marvel/mclient/operations"
	"github.com/loivis/mcapi-loader/marvel/models"
	"github.com/loivis/mcapi-loader/mcapiloader"
)

func (p *Processor) loadCharacters(ctx context.Context) error {
	err := p.loadAllCharactersWithBasicInfo(ctx)
	if err != nil {
		return fmt.Errorf("error loading all characters info: %v", err)
	}

	return nil
}

func (p *Processor) loadAllCharactersWithBasicInfo(ctx context.Context) error {
	remote, err := p.getCharacterCount(ctx)
	if err != nil {
		return fmt.Errorf("error fetching character count: %v", err)
	}
	log.Info().Str("type", "character").Int32("count", remote).Msg("fetched")

	existing, err := p.store.GetCount("characters")
	if err != nil {
		return err
	}
	log.Info().Str("type", "character").Int64("count", existing).Msg("existing characters")

	if int64(remote) == existing {
		log.Info().Int64("local", existing).Int32("remote", remote).Msg("no missing characters")
		return nil
	}

	log.Info().Int64("local", existing).Int32("remote", remote).Msg("missing characters, reload")
	chars, err := p.getAllCharacters(ctx, remote)
	if err != nil {
		return fmt.Errorf("error getting all characters: %v", err)
	}

	log.Info().Int("count", len(chars)).Msg("fetched all characters")

	err = p.store.SaveCharacters(chars)
	if err != nil {
		return fmt.Errorf("error storing characters: %v", err)
	}

	return nil
}

func (p *Processor) getCharacterCount(ctx context.Context) (int32, error) {
	var limit int32 = 1
	params := &operations.GetCharactersCollectionParams{
		Limit: &limit,
	}
	p.setParams(ctx, params)

	col, err := p.mclient.Operations.GetCharactersCollection(params)
	if err != nil {
		return 0, err
	}

	return col.Payload.Data.Total, nil
}

func (p *Processor) getAllCharacters(ctx context.Context, count int32) ([]*mcapiloader.Character, error) {
	chars := make([]*mcapiloader.Character, count)

	var g errgroup.Group

	for i := 0; i < int(count/p.limit)+1; i++ {
		offset := p.limit * int32(i)
		g.Go(func() error {
			paged, err := p.getPagedCharacters(ctx, offset)
			if err != nil {
				return err
			}

			for j, char := range paged {
				chars[offset+int32(j)] = char
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return chars, nil
}

func (p *Processor) getPagedCharacters(ctx context.Context, offset int32) ([]*mcapiloader.Character, error) {
	chars := []*mcapiloader.Character{}

	params := &operations.GetCharactersCollectionParams{
		Limit:  &p.limit,
		Offset: &offset,
	}
	p.setParams(ctx, params)

	col, err := p.mclient.Operations.GetCharactersCollection(params)
	if err != nil {
		return nil, err
	}

	for _, res := range col.Payload.Data.Results {
		char, err := convertCharacter(res)
		if err != nil {
			return nil, fmt.Errorf("error converting character %s(%d): %v", res.Name, res.ID, err)
		}
		chars = append(chars, char)
	}

	log.Info().Int32("offset", offset).Int("count", len(chars)).Msg("fetched")

	return chars, nil
}

func convertCharacter(in *models.Character) (*mcapiloader.Character, error) {
	out := &mcapiloader.Character{
		Description: in.Description,
		ID:          in.ID,
		Modified:    in.Modified,
		Name:        in.Name,
		Thumbnail:   in.Thumbnail.Path + "." + in.Thumbnail.Extension,
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
		out.URLs = append(out.URLs, &mcapiloader.URL{
			Type: url.Type,
			URL:  strings.Split(url.URL, "?")[0],
		})
	}

	if in.Comics.Available == in.Comics.Returned && in.Events.Available == in.Events.Returned && in.Series.Available == in.Series.Returned && in.Stories.Available == in.Stories.Returned {
		out.Intact = true
	}

	return out, nil
}
