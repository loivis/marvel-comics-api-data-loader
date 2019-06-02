package marvel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		c := NewClient("http://base.url", "private", "public")

		if got, want := c.baseURL, "http://base.url"; got != want {
			t.Errorf("got baseURL %q, want %q", got, want)
		}

		if got, want := c.privateKey, "private"; got != want {
			t.Errorf("got privateKey %q, want %q", got, want)
		}

		if got, want := c.publicKey, "public"; got != want {
			t.Errorf("got publicKey %q, want %q", got, want)
		}
	})

	t.Run("TrimTrailingSlash", func(t *testing.T) {
		c := NewClient("http://base.url//", "private", "public")

		if got, want := c.baseURL, "http://base.url"; got != want {
			t.Errorf("got baseURL %q, want %q", got, want)
		}

		if got, want := c.privateKey, "private"; got != want {
			t.Errorf("got privateKey %q, want %q", got, want)
		}

		if got, want := c.publicKey, "public"; got != want {
			t.Errorf("got publicKey %q, want %q", got, want)
		}
	})
}

func TestClient_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var gotRequest *http.Request

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotRequest = r
			w.Write([]byte(`{"foo":"bar"}`))
		}))

		c := NewClient(ts.URL, "", "public")

		b, err := c.get(context.Background(), &Params{typ: "foo", Limit: 1, Offset: 2})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if got, want := gotRequest.Method, http.MethodGet; got != want {
			t.Errorf("got request method %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Path, "/foo"; got != want {
			t.Errorf("got request method %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("apikey"), c.publicKey; got != want {
			t.Errorf("got apikey %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("limit"), "1"; got != want {
			t.Errorf("got limit %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("offset"), "2"; got != want {
			t.Errorf("got offset %q, want %q", got, want)
		}

		if got, want := string(b), string(b); got != want {
			t.Errorf("got bytes %q, want %q", got, want)
		}
	})

	t.Run("APIError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusBadRequest)
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.get(context.Background(), &Params{typ: "foo"})

		if err == nil {
			t.Fatal("error is nil")
		}

		ae, ok := err.(*APIError)

		if !ok {
			t.Fatalf("got error type %T, want %T", ae, err)
		}

		if got, want := ae.Code, http.StatusBadRequest; got != want {
			t.Errorf("got status code %d, want %d", got, want)
		}
	})

	t.Run("PathError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusBadRequest)
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.get(context.Background(), nil)

		if err == nil {
			t.Fatal("error is nil")
		}

		pe, ok := err.(*PathError)

		if !ok {
			t.Fatalf("got error type %T, want %T", pe, err)
		}

		if got, want := pe.Error(), "missing in path: type, id, subtype"; got != want {
			t.Errorf("got error %q, want %q", got, want)
		}
	})
}

func TestClient_GetCount(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var gotRequest *http.Request

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotRequest = r
			w.Write([]byte(`{"data":{"total":123}}`))
		}))

		c := NewClient(ts.URL, "", "public")

		count, err := c.GetCount(context.Background(), "foo")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if got, want := gotRequest.Method, http.MethodGet; got != want {
			t.Errorf("got request method %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Path, "/foo"; got != want {
			t.Errorf("got request method %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("apikey"), c.publicKey; got != want {
			t.Errorf("got apikey %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("limit"), "1"; got != want {
			t.Errorf("got limit %q, want %q", got, want)
		}

		if got, want := count, 123; got != want {
			t.Errorf("got count %d, want %d", got, want)
		}
	})

	t.Run("APIError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusBadRequest)
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.GetCount(context.Background(), "foo")

		if err == nil {
			t.Fatal("error is nil")
		}

		ae, ok := err.(*APIError)

		if !ok {
			t.Fatalf("got error type %T, want %T", ae, err)
		}

		if got, want := ae.Code, http.StatusBadRequest; got != want {
			t.Errorf("got status code %d, want %d", got, want)
		}
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{invalid-json`))
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.GetCount(context.Background(), "foo")

		if err == nil {
			t.Fatal("error is nil")
		}

		if err == nil {
			t.Fatal("error is nil")
		}

		if v, ok := err.(*json.SyntaxError); !ok {
			t.Fatalf("got error type %T, want %T", v, err)
		}
	})
}

func TestClient_BuildQuery(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		params *Params
		want   url.Values
	}{
		{
			desc:   "Success",
			params: &Params{Limit: 1, Offset: 2},
			want: url.Values{
				"apikey": {"public"},
				"foo":    {"bar", "baz"},
				"hash":   {"fc1af084cf4c71dd2d6ddbd14e5b3e32"},
				"limit":  {"1"},
				"offset": {"2"},
				"ts":     {"12345"},
			},
		},
		{
			desc:   "NilParams",
			params: nil,
			want: url.Values{
				"apikey": {"public"},
				"foo":    {"bar", "baz"},
				"hash":   {"fc1af084cf4c71dd2d6ddbd14e5b3e32"},
				"ts":     {"12345"},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			c := NewClient("", "", "public")

			got := c.buildQuery(url.Values{"foo": {"bar", "baz"}}, tc.params, "12345")

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got query %v, want %v", got, tc.want)
			}
		})
	}
}

func Test_PathFromParams(t *testing.T) {
	id := 123

	t.Run("Success", func(t *testing.T) {
		for _, tc := range []struct {
			desc  string
			paths *Params
			want  string
		}{
			{
				desc:  "TypeOnly",
				paths: &Params{typ: "type"},
				want:  "type",
			},
			{
				desc:  "TypeAndID",
				paths: &Params{typ: "type", id: &id},
				want:  "type/123",
			},
			{
				desc:  "TypeAndIDAndSubType",
				paths: &Params{typ: "type", id: &id, subtype: "sub-type"},
				want:  "type/123/sub-type",
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				got, err := pathFromParams(tc.paths)
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}

				if got, want := got, tc.want; got != want {
					t.Errorf("got path %q, want %q", got, want)
				}
			})
		}
	})

	t.Run("Error", func(t *testing.T) {
		for _, tc := range []struct {
			desc  string
			paths *Params
		}{
			{
				desc:  "Nil",
				paths: nil,
			},
			{
				desc:  "MissingTypeID",
				paths: &Params{subtype: "sub-type"},
			},
			{
				desc:  "MissingType",
				paths: &Params{id: &id, subtype: "sub-type"},
			},
			{
				desc:  "MissingID",
				paths: &Params{typ: "type", subtype: "sub-type"},
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				_, err := pathFromParams(tc.paths)
				if err == nil {
					t.Fatal("error is nil")
				}
			})
		}
	})
}
