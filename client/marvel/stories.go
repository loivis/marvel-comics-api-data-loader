package marvel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
)

// GetStory returns the story of specified id with given Params.
func (c *Client) GetStory(ctx context.Context, id int) (*Story, error) {
	resp, err := c.get(ctx, &Params{typ: m27r.TypeStories, id: &id})
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Story
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Data.Results) == 0 {
		return nil, errors.New("no story returned")
	}

	return data.Data.Results[0], nil
}

// GetStories returns list of story with given Params.
func (c *Client) GetStories(ctx context.Context, params *Params) ([]*Story, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeStories

	return c.getStories(ctx, params)
}

// GetCharacterStories returns list of story filtered by character ID.
func (c *Client) GetCharacterStories(ctx context.Context, id int, params *Params) ([]*Story, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeCharacters
	params.id = &id
	params.subtype = m27r.TypeStories

	return c.getStories(ctx, params)
}

// GetComicStories returns list of story filtered by comic ID.
func (c *Client) GetComicStories(ctx context.Context, id int, params *Params) ([]*Story, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeComics
	params.id = &id
	params.subtype = m27r.TypeStories

	return c.getStories(ctx, params)
}

// GetCreatorStories returns list of story filtered by creator ID.
func (c *Client) GetCreatorStories(ctx context.Context, id int, params *Params) ([]*Story, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeCreators
	params.id = &id
	params.subtype = m27r.TypeStories

	return c.getStories(ctx, params)
}

// GetEventStories returns list of story filtered by event ID.
func (c *Client) GetEventStories(ctx context.Context, id int, params *Params) ([]*Story, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeEvents
	params.id = &id
	params.subtype = m27r.TypeStories

	return c.getStories(ctx, params)
}

// GetSeriesStories returns list of story filtered by series ID.
func (c *Client) GetSeriesStories(ctx context.Context, id int, params *Params) ([]*Story, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeSeries
	params.id = &id
	params.subtype = m27r.TypeStories

	return c.getStories(ctx, params)
}

// getStories returns a list of story with given Params.
// If Params.subTyp is set, it returns stories filtered by Params.Type and Params.ID.
func (c *Client) getStories(ctx context.Context, params *Params) ([]*Story, error) {
	resp, err := c.get(ctx, params)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Story
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data.Data.Results, nil
}
