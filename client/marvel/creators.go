package marvel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
)

// GetCreator returns the creator of specified id with given Params.
func (c *Client) GetCreator(ctx context.Context, id int) (*Creator, error) {
	resp, err := c.get(ctx, &Params{typ: m27r.TypeCreators, id: &id})
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Creator
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Data.Results) == 0 {
		return nil, errors.New("no creator returned")
	}

	return data.Data.Results[0], nil
}

// GetCreators returns list of creator with given Params.
func (c *Client) GetCreators(ctx context.Context, params *Params) ([]*Creator, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeCreators

	return c.getCreators(ctx, params)
}

// GetComicCreators returns list of creators filtered by comic ID.
func (c *Client) GetComicCreators(ctx context.Context, id int, params *Params) ([]*Creator, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeComics
	params.id = &id
	params.subtype = m27r.TypeCreators

	return c.getCreators(ctx, params)
}

// GetEventCreators returns list of creators filtered by event ID.
func (c *Client) GetEventCreators(ctx context.Context, id int, params *Params) ([]*Creator, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeEvents
	params.id = &id
	params.subtype = m27r.TypeCreators

	return c.getCreators(ctx, params)
}

// GetSeriesCreators returns list of creators filtered by series ID.
func (c *Client) GetSeriesCreators(ctx context.Context, id int, params *Params) ([]*Creator, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeSeries
	params.id = &id
	params.subtype = m27r.TypeCreators

	return c.getCreators(ctx, params)
}

// GetStoryCreators returns list of creators filtered by story ID.
func (c *Client) GetStoryCreators(ctx context.Context, id int, params *Params) ([]*Creator, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeStories
	params.id = &id
	params.subtype = m27r.TypeCreators

	return c.getCreators(ctx, params)
}

// getCreators returns a list of creator with given Params.
// If Params.subTyp is set, it returns creators filtered by Params.Type and Params.ID.
func (c *Client) getCreators(ctx context.Context, params *Params) ([]*Creator, error) {
	resp, err := c.get(ctx, params)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Creator
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data.Data.Results, nil
}
