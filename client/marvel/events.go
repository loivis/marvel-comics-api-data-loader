package marvel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
)

// GetEvent returns the event of specified id with given Params.
func (c *Client) GetEvent(ctx context.Context, id int) (*Event, error) {
	resp, err := c.get(ctx, &Params{typ: m27r.TypeEvents, id: &id})
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Event
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	if len(data.Data.Results) == 0 {
		return nil, errors.New("no event returned")
	}

	return data.Data.Results[0], nil
}

// GetEvents returns list of event with given Params.
func (c *Client) GetEvents(ctx context.Context, params *Params) ([]*Event, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeEvents

	return c.getEvents(ctx, params)
}

// GetCharacterEvents returns list of events filtered by character ID.
func (c *Client) GetCharacterEvents(ctx context.Context, id int, params *Params) ([]*Event, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeCharacters
	params.id = &id
	params.subtype = m27r.TypeEvents

	return c.getEvents(ctx, params)
}

// GetComicEvents returns list of events filtered by comic ID.
func (c *Client) GetComicEvents(ctx context.Context, id int, params *Params) ([]*Event, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeComics
	params.id = &id
	params.subtype = m27r.TypeEvents

	return c.getEvents(ctx, params)
}

// GetCreatorEvents returns list of events filtered by creator ID.
func (c *Client) GetCreatorEvents(ctx context.Context, id int, params *Params) ([]*Event, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeCreators
	params.id = &id
	params.subtype = m27r.TypeEvents

	return c.getEvents(ctx, params)
}

// GetSeriesEvents returns list of events filtered by series ID.
func (c *Client) GetSeriesEvents(ctx context.Context, id int, params *Params) ([]*Event, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeSeries
	params.id = &id
	params.subtype = m27r.TypeEvents

	return c.getEvents(ctx, params)
}

// GetStoryEvents returns list of events filtered by story ID.
func (c *Client) GetStoryEvents(ctx context.Context, id int, params *Params) ([]*Event, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	params.typ = m27r.TypeStories
	params.id = &id
	params.subtype = m27r.TypeEvents

	return c.getEvents(ctx, params)
}

// getEvents returns a list of event with given Params.
// If Params.subTyp is set, it returns events filtered by Params.Type and Params.ID.
func (c *Client) getEvents(ctx context.Context, params *Params) ([]*Event, error) {
	resp, err := c.get(ctx, params)
	if err != nil {
		return nil, err
	}

	var data struct {
		Data struct {
			Results []*Event
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}

	return data.Data.Results, nil
}
