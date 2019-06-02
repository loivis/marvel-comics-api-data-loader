package marvel

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestComic_UnmarshalJSON(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		for _, tc := range []struct {
			desc string
			b    []byte
			want Comic
		}{
			{
				desc: "Nil",
				b:    []byte(`{"id":1}`),
				want: Comic{ID: 1},
			},
			{
				desc: "String",
				b:    []byte(`{"id":1,"diamondCode":"foo", "isbn":"bar"}`),
				want: Comic{ID: 1, DiamondCode: "foo", ISBN: "bar"},
			},
			{
				desc: "Int",
				b:    []byte(`{"id":1,"diamondCode":1234, "isbn":5678}`),
				want: Comic{ID: 1, DiamondCode: "1234", ISBN: "5678"},
			},
			{
				desc: "Float",
				b:    []byte(`{"id":1,"diamondCode":12.34, "isbn":56.78}`),
				want: Comic{ID: 1, DiamondCode: "12.34", ISBN: "56.78"},
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				var gotComic Comic

				err := json.Unmarshal(tc.b, &gotComic)
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}

				if got, want := gotComic.ID, tc.want.ID; got != want {
					t.Errorf("got ID %d, want %d", got, want)
				}

				if got, want := gotComic.DiamondCode, tc.want.DiamondCode; got != want {
					t.Errorf("got DiamondCode %q, want %q", got, want)
				}

				if got, want := gotComic.ISBN, tc.want.ISBN; got != want {
					t.Errorf("got ISBN %q, want %q", got, want)
				}

			})
		}
	})

	t.Run("Error", func(t *testing.T) {
		for _, tc := range []struct {
			desc string
			in   []byte
			err  error
		}{
			{
				desc: "MalformedJSON",
				in:   []byte(`{invalid-json`),
				err:  &json.SyntaxError{},
			},
			{
				desc: "InvalidTypeID",
				in:   []byte(`{"id":false}`),
				err:  &json.UnmarshalTypeError{},
			},
			{
				desc: "InvalidTypeDiamondCode",
				in:   []byte(`{"id":1, "DiamondCode":false}`),
				err:  &json.UnmarshalTypeError{},
			},
			{
				desc: "InvalidTypeISBN",
				in:   []byte(`{"id":1, "ISBN":false}`),
				err:  &json.UnmarshalTypeError{},
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				var comic Comic

				gotErr := json.Unmarshal(tc.in, &comic)

				if gotErr == nil {
					t.Fatal("error is nil")
				}

				if got, want := reflect.TypeOf(gotErr).String(), reflect.TypeOf(tc.err).String(); got != want {
					t.Errorf("got error type %q, want %q", got, want)
				}
			})
		}
	})
}

func TestStory_UnmarshalJSON(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		for _, tc := range []struct {
			desc string
			b    []byte
			want Story
		}{
			{
				desc: "Nil",
				b:    []byte(`{"id":1}`),
				want: Story{ID: 1},
			},
			{
				desc: "String",
				b:    []byte(`{"id":1,"title":"foo"}`),
				want: Story{ID: 1, Title: "foo"},
			},
			{
				desc: "Int",
				b:    []byte(`{"id":1,"title":1234}`),
				want: Story{ID: 1, Title: "1234"},
			},
			{
				desc: "Float",
				b:    []byte(`{"id":1,"title":12.34}`),
				want: Story{ID: 1, Title: "12.34"},
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				var gotStory Story

				err := json.Unmarshal(tc.b, &gotStory)
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}

				if got, want := gotStory.ID, tc.want.ID; got != want {
					t.Errorf("got ID %d, want %d", got, want)
				}

				if got, want := gotStory.Title, tc.want.Title; got != want {
					t.Errorf("got Title %q, want %q", got, want)
				}
			})
		}
	})

	t.Run("Error", func(t *testing.T) {
		for _, tc := range []struct {
			desc string
			in   []byte
			err  error
		}{
			{
				desc: "MalformedJSON",
				in:   []byte(`{invalid-json`),
				err:  &json.SyntaxError{},
			},
			{
				desc: "InvalidTypeID",
				in:   []byte(`{"id":false}`),
				err:  &json.UnmarshalTypeError{},
			},
			{
				desc: "InvalidTypeTitle",
				in:   []byte(`{"id":1, "Title":false}`),
				err:  &json.UnmarshalTypeError{},
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				var story Story

				gotErr := json.Unmarshal(tc.in, &story)

				if gotErr == nil {
					t.Fatal("error is nil")
				}

				if got, want := reflect.TypeOf(gotErr).String(), reflect.TypeOf(tc.err).String(); got != want {
					t.Errorf("got error type %q, want %q", got, want)
				}
			})
		}
	})
}

func TestCreator_UnmarshalJSON(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		for _, tc := range []struct {
			desc string
			b    []byte
			want Creator
		}{
			{
				desc: "Nil",
				b:    []byte(`{"id":1}`),
				want: Creator{ID: 1},
			},
			{
				desc: "String",
				b:    []byte(`{"id":1,"lastName":"foo", "suffix":"bar"}`),
				want: Creator{ID: 1, LastName: "foo", Suffix: "bar"},
			},
			{
				desc: "Int",
				b:    []byte(`{"id":1,"lastName":1234, "suffix":5678}`),
				want: Creator{ID: 1, LastName: "1234", Suffix: "5678"},
			},
			{
				desc: "Float",
				b:    []byte(`{"id":1,"lastName":12.34, "suffix":56.78}`),
				want: Creator{ID: 1, LastName: "12.34", Suffix: "56.78"},
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				var gotCreator Creator

				err := json.Unmarshal(tc.b, &gotCreator)
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}

				if got, want := gotCreator.ID, tc.want.ID; got != want {
					t.Errorf("got ID %d, want %d", got, want)
				}

				if got, want := gotCreator.LastName, tc.want.LastName; got != want {
					t.Errorf("got LastName %q, want %q", got, want)
				}

				if got, want := gotCreator.Suffix, tc.want.Suffix; got != want {
					t.Errorf("got Suffix %q, want %q", got, want)
				}

			})
		}
	})

	t.Run("Error", func(t *testing.T) {
		for _, tc := range []struct {
			desc string
			in   []byte
			err  error
		}{
			{
				desc: "MalformedJSON",
				in:   []byte(`{invalid-json`),
				err:  &json.SyntaxError{},
			},
			{
				desc: "InvalidTypeID",
				in:   []byte(`{"id":false}`),
				err:  &json.UnmarshalTypeError{},
			},
			{
				desc: "InvalidTypeLastName",
				in:   []byte(`{"id":1, "LastName":false}`),
				err:  &json.UnmarshalTypeError{},
			},
			{
				desc: "InvalidTypeSuffix",
				in:   []byte(`{"id":1, "Suffix":false}`),
				err:  &json.UnmarshalTypeError{},
			},
		} {
			t.Run(tc.desc, func(t *testing.T) {
				var creator Creator

				gotErr := json.Unmarshal(tc.in, &creator)

				if gotErr == nil {
					t.Fatal("error is nil")
				}

				if got, want := reflect.TypeOf(gotErr).String(), reflect.TypeOf(tc.err).String(); got != want {
					t.Errorf("got error type %q, want %q", got, want)
				}
			})
		}
	})
}
