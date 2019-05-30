package main

import (
	"strings"
	"testing"
)

func Test_Current(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := strings.NewReader(`{"foo":{"offset":1,"ids":[1,2,3]}}`)
		gotMappings, err := current(r)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wantMappings := map[string]*collection{
			"foo": &collection{
				Offset: 1,
				IDs:    []int{1, 2, 3},
			},
			"characters": &collection{},
			"comics":     &collection{},
			"creators":   &collection{},
			"events":     &collection{},
			"series":     &collection{},
			"stories":    &collection{},
		}

		if got, want := len(gotMappings), len(wantMappings); got != want {
			t.Fatalf("got %d mappings, want %d", got, want)
		}

		for k, col := range wantMappings {
			if got, want := gotMappings[k].Offset, col.Offset; got != want {
				t.Errorf("got mappings[%s].Offset %d, want %d", k, got, want)
			}

			if got, want := len(gotMappings[k].IDs), len(col.IDs); got != want {
				t.Errorf("mappings[%s].IDs has %d elements, want %d", k, got, want)
			}

			for i, id := range col.IDs {
				if got, want := gotMappings[k].IDs[i], id; got != want {
					t.Errorf("got mappins[%s].IDs[%d] %d, want %d", k, i, got, want)
				}
			}
		}
	})

	t.Run("MalformedJSON", func(t *testing.T) {
		r := strings.NewReader(`{`)
		_, err := current(r)

		if err == nil {
			t.Fatal("error is nil")
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		r := strings.NewReader(``)
		mappings, err := current(r)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if mappings == nil {
			t.Fatal("mappings is nil")
		}
		if got, want := len(mappings), 6; got != want {
			t.Errorf("got %d mappings, want %d", got, want)
		}
	})
}
