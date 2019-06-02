package marvel

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

// Client contains info for requests to marvel comics api.
type Client struct {
	baseURL    string
	hc         *http.Client
	privateKey string
	publicKey  string
}

// Params contains path info and query parameters for api requests.
type Params struct {
	typ     string // private field
	id      *int   // avoid 0 when not set
	subtype string

	Limit   int
	Offset  int
	OrderBy string
}

// NewClient returns a marvel Client.
func NewClient(baseURL, privateKey, publicKey string) *Client {
	return &Client{
		baseURL:    strings.Trim(baseURL, "/"),
		hc:         &http.Client{Timeout: 150 * time.Second},
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// GetCount returns count for the given type.
func (c *Client) GetCount(ctx context.Context, typ string) (int, error) {
	resp, err := c.get(ctx, &Params{typ: typ, Limit: 1})
	if err != nil {
		return 0, err
	}

	var data struct {
		Data struct {
			Total int
		}
	}

	err = json.Unmarshal(resp, &data)
	if err != nil {
		return 0, err
	}

	return data.Data.Total, nil
}

// get returns api response as slice of byte with given path and request params.
func (c *Client) get(ctx context.Context, params *Params) ([]byte, error) {
	path, err := pathFromParams(params)
	if err != nil {
		return nil, err
	}

	u := fmt.Sprintf("%s/%s", c.baseURL, path)

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	ts := fmt.Sprintf("%d", time.Now().Unix())
	req.URL.RawQuery = c.buildQuery(req.URL.Query(), params, ts).Encode()

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &APIError{Code: resp.StatusCode, Message: string(b)}
	}

	return b, nil
}

func (c *Client) buildQuery(in url.Values, params *Params, ts string) url.Values {
	out := url.Values{
		"apikey": {c.publicKey},
		"hash":   {fmt.Sprintf("%x", md5.Sum([]byte(ts+c.privateKey+c.publicKey)))},
		"ts":     {ts},
	}

	for k, v := range in {
		out[k] = append(out[k], v...)
	}

	if params != nil {
		if params.Limit != 0 {
			out["limit"] = []string{strconv.FormatInt(int64(params.Limit), 10)}
		}

		if params.Offset != 0 {
			out["offset"] = []string{strconv.FormatInt(int64(params.Offset), 10)}
		}

		if params.OrderBy != "" {
			out["orderBy"] = []string{params.OrderBy}
		}
	}

	return out
}

func pathFromParams(params *Params) (string, error) {
	switch {
	case params == nil:
		return "", &PathError{Missing: []string{"type", "id", "subtype"}}
	case params.typ == "" && params.id == nil && params.subtype != "":
		return "", &PathError{Missing: []string{"type", "id"}}
	case params.typ == "":
		return "", &PathError{Missing: []string{"type"}}
	case params.typ != "" && params.id == nil && params.subtype != "":
		return "", &PathError{Missing: []string{"id"}}
	}

	if params.id == nil {
		return params.typ, nil
	}

	return path.Join(params.typ, strconv.FormatInt(int64(*params.id), 10), params.subtype), nil
}
