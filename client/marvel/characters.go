package marvel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
)

// GetCharacter returns the character of specified id with given Params.
func (c *Client) GetCharacter(ctx context.Context, id int) (*Character, error) {
	resp, err := c.get(ctx, &Params{typ: m27r.TypeCharacters, id: &id})
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Character
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Data.Results) == 0 {
		return nil, errors.New("no character returned")
	}

	return data.Data.Results[0], nil
}

// GetCharacters returns list of character with given Params.
func (c *Client) GetCharacters(ctx context.Context, params *Params) ([]*Character, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeCharacters

	return c.getCharacters(ctx, params)
}

// GetComicCharacters returns list of character filtered by comic ID.
func (c *Client) GetComicCharacters(ctx context.Context, id int, params *Params) ([]*Character, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeComics
	params.id = &id
	params.subtype = m27r.TypeCharacters

	return c.getCharacters(ctx, params)
}

// GetEventCharacters returns list of character filtered by event ID.
func (c *Client) GetEventCharacters(ctx context.Context, id int, params *Params) ([]*Character, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeEvents
	params.id = &id
	params.subtype = m27r.TypeCharacters

	return c.getCharacters(ctx, params)
}

// GetSeriesCharacters returns list of character filtered by series ID.
func (c *Client) GetSeriesCharacters(ctx context.Context, id int, params *Params) ([]*Character, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeSeries
	params.id = &id
	params.subtype = m27r.TypeCharacters

	return c.getCharacters(ctx, params)
}

// GetStoryCharacters returns list of character filtered by story ID.
func (c *Client) GetStoryCharacters(ctx context.Context, id int, params *Params) ([]*Character, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeStories
	params.id = &id
	params.subtype = m27r.TypeCharacters

	return c.getCharacters(ctx, params)
}

// getCharacters returns a list of character with given Params.
// If Params.subTyp is set, it returns characters filtered by Params.Type and Params.ID.
func (c *Client) getCharacters(ctx context.Context, params *Params) ([]*Character, error) {
	resp, err := c.get(ctx, params)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Character
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data.Data.Results, nil
}
