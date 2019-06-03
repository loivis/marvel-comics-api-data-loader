package marvel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/loivis/marvel-comics-api-data-loader/maco"
)

// GetComic returns the comic of specified id with given Params.
func (c *Client) GetComic(ctx context.Context, id int) (*Comic, error) {
	resp, err := c.get(ctx, &Params{typ: maco.TypeComics, id: &id})
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Comic
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Data.Results) == 0 {
		return nil, errors.New("no comic returned")
	}

	return data.Data.Results[0], nil
}

// GetComics returns list of comic with given Params.
func (c *Client) GetComics(ctx context.Context, params *Params) ([]*Comic, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeComics

	return c.getComics(ctx, params)
}

// GetCharacterComics returns list of comics filtered by character ID.
func (c *Client) GetCharacterComics(ctx context.Context, id int, params *Params) ([]*Comic, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeCharacters
	params.id = &id
	params.subtype = maco.TypeComics

	return c.getComics(ctx, params)
}

// GetCreatorComics returns list of comics filtered by creator ID.
func (c *Client) GetCreatorComics(ctx context.Context, id int, params *Params) ([]*Comic, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeCreators
	params.id = &id
	params.subtype = maco.TypeComics

	return c.getComics(ctx, params)
}

// GetEventComics returns list of comics filtered by event ID.
func (c *Client) GetEventComics(ctx context.Context, id int, params *Params) ([]*Comic, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeEvents
	params.id = &id
	params.subtype = maco.TypeComics

	return c.getComics(ctx, params)
}

// GetSeriesComics returns list of comics filtered by series ID.
func (c *Client) GetSeriesComics(ctx context.Context, id int, params *Params) ([]*Comic, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeSeries
	params.id = &id
	params.subtype = maco.TypeComics

	return c.getComics(ctx, params)
}

// GetStoryComics returns list of comics filtered by story ID.
func (c *Client) GetStoryComics(ctx context.Context, id int, params *Params) ([]*Comic, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeStories
	params.id = &id
	params.subtype = maco.TypeComics

	return c.getComics(ctx, params)
}

// getComics returns a list of comic with given Params.
// If Params.subTyp is set, it returns comics filtered by Params.Type and Params.ID.
func (c *Client) getComics(ctx context.Context, params *Params) ([]*Comic, error) {
	resp, err := c.get(ctx, params)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Comic
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data.Data.Results, nil
}
