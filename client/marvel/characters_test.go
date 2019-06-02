package marvel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_GetCharacter(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var gotRequest *http.Request

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotRequest = r
			w.Write([]byte(`{"data":{"results":[{"id":1}]}}`))
		}))

		c := NewClient(ts.URL, "", "")

		gotChar, err := c.GetCharacter(context.Background(), 123)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		wantChar := &Character{ID: 1}

		if got, want := gotRequest.URL.Path, "/characters/123"; got != want {
			t.Errorf("got request path %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("apikey"), c.publicKey; got != want {
			t.Errorf("got apikey %q, want %q", got, want)
		}

		if got, want := gotChar.ID, wantChar.ID; got != want {
			t.Errorf("got character ID %d, want %d", got, want)
		}
	})

	t.Run("NoReturnError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"data":{"results":[]}}`))
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.GetCharacter(context.Background(), 123)

		if err == nil {
			t.Fatal("error is nil")
		}
	})

	t.Run("APIError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusBadRequest)
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.GetCharacter(context.Background(), 123)

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

		_, err := c.GetCharacter(context.Background(), 123)

		if err == nil {
			t.Fatal("error is nil")
		}

		if v, ok := err.(*json.SyntaxError); !ok {
			t.Fatalf("got error type %T, want %T", v, err)
		}
	})
}

func TestClient_GetCharacters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var gotRequest *http.Request

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotRequest = r
			w.Write([]byte(`{"data":{"results":[{"id":1},{"id":2}]}}`))
		}))

		c := NewClient(ts.URL, "", "")

		gotChars, err := c.GetCharacters(context.Background(), &Params{Limit: 1, Offset: 2})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		wantChars := []*Character{
			&Character{ID: 1},
			&Character{ID: 2},
		}

		if got, want := gotRequest.URL.Path, "/characters"; got != want {
			t.Errorf("got request path %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("apikey"), c.publicKey; got != want {
			t.Errorf("got apikey %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("limit"), "1"; got != want {
			t.Errorf("got query parameter limit %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("offset"), "2"; got != want {
			t.Errorf("got query parameter offset %q, want %q", got, want)
		}

		if got, want := len(gotChars), len(wantChars); got != want {
			t.Errorf("got %d characters, want %d", got, want)
		}

		for i, char := range wantChars {
			if got, want := gotChars[i].ID, char.ID; got != want {
				t.Errorf("got chars[%d].ID %d, want %d", i, got, want)
			}
		}
	})

	t.Run("NilParamsError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusBadRequest)
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.GetCharacters(context.Background(), nil)

		if err == nil {
			t.Fatal("error is nil")
		}
	})

	t.Run("APIError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusBadRequest)
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.GetCharacters(context.Background(), &Params{})

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

		_, err := c.GetCharacters(context.Background(), &Params{})

		if err == nil {
			t.Fatal("error is nil")
		}

		if v, ok := err.(*json.SyntaxError); !ok {
			t.Fatalf("got error type %T, want %T", v, err)
		}
	})
}

func TestClient_GetComicCharacters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var gotRequest *http.Request

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotRequest = r
			w.Write([]byte(`{"data":{"results":[{"id":1},{"id":2}]}}`))
		}))

		c := NewClient(ts.URL, "", "")

		gotChars, err := c.GetComicCharacters(context.Background(), 123, &Params{Limit: 1, Offset: 2})
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		wantChars := []*Character{
			&Character{ID: 1},
			&Character{ID: 2},
		}

		if got, want := gotRequest.URL.Path, "/comics/123/characters"; got != want {
			t.Errorf("got request path %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("apikey"), c.publicKey; got != want {
			t.Errorf("got apikey %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("limit"), "1"; got != want {
			t.Errorf("got query parameter limit %q, want %q", got, want)
		}

		if got, want := gotRequest.URL.Query().Get("offset"), "2"; got != want {
			t.Errorf("got query parameter offset %q, want %q", got, want)
		}

		if got, want := len(gotChars), len(wantChars); got != want {
			t.Errorf("got %d characters, want %d", got, want)
		}

		for i, char := range wantChars {
			if got, want := gotChars[i].ID, char.ID; got != want {
				t.Errorf("got chars[%d].ID %d, want %d", i, got, want)
			}
		}
	})

	t.Run("NilParamsError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusBadRequest)
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.GetComicCharacters(context.Background(), 123, nil)

		if err == nil {
			t.Fatal("error is nil")
		}
	})

	t.Run("APIError", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "", http.StatusBadRequest)
		}))

		c := NewClient(ts.URL, "", "")

		_, err := c.GetComicCharacters(context.Background(), 123, &Params{})

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

		_, err := c.GetComicCharacters(context.Background(), 123, &Params{})

		if err == nil {
			t.Fatal("error is nil")
		}

		if v, ok := err.(*json.SyntaxError); !ok {
			t.Fatalf("got error type %T, want %T", v, err)
		}
	})
}
