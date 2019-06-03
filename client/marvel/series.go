package marvel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/loivis/marvel-comics-api-data-loader/maco"
)

// GetSeries returns the series of specified id with given Params.
func (c *Client) GetSeriesSingle(ctx context.Context, id int) (*Series, error) {
	resp, err := c.get(ctx, &Params{typ: maco.TypeSeries, id: &id})
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Series
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Data.Results) == 0 {
		return nil, errors.New("no series returned")
	}

	return data.Data.Results[0], nil
}

// GetSeries returns list of series with given Params.
func (c *Client) GetSeries(ctx context.Context, params *Params) ([]*Series, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeSeries

	return c.getSeries(ctx, params)
}

// GetCharacterSeries returns list of series filtered by character ID.
func (c *Client) GetCharacterSeries(ctx context.Context, id int, params *Params) ([]*Series, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeCharacters
	params.id = &id
	params.subtype = maco.TypeSeries

	return c.getSeries(ctx, params)
}

// GetComicSeries returns list of series filtered by comic ID.
func (c *Client) GetComicSeries(ctx context.Context, id int, params *Params) ([]*Series, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeComics
	params.id = &id
	params.subtype = maco.TypeSeries

	return c.getSeries(ctx, params)
}

// GetCreatorSeries returns list of series filtered by creator ID.
func (c *Client) GetCreatorSeries(ctx context.Context, id int, params *Params) ([]*Series, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeCreators
	params.id = &id
	params.subtype = maco.TypeSeries

	return c.getSeries(ctx, params)
}

// GetEventSeries returns list of series filtered by event ID.
func (c *Client) GetEventSeries(ctx context.Context, id int, params *Params) ([]*Series, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeEvents
	params.id = &id
	params.subtype = maco.TypeSeries

	return c.getSeries(ctx, params)
}

// GetStorySeries returns list of series filtered by story ID.
func (c *Client) GetStorySeries(ctx context.Context, id int, params *Params) ([]*Series, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = maco.TypeStories
	params.id = &id
	params.subtype = maco.TypeSeries

	return c.getSeries(ctx, params)
}

// getSeries returns a list of series with given Params.
// If Params.subTyp is set, it returns seriess filtered by Params.Type and Params.ID.
func (c *Client) getSeries(ctx context.Context, params *Params) ([]*Series, error) {
	resp, err := c.get(ctx, params)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Series
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data.Data.Results, nil
}
