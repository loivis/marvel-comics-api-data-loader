package process

import (
	"testing"
)

func Test_IDFromURL(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id, err := idFromURL("https://example.com/123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := id, 123; got != want {
			t.Errorf("got id %d, want %d", got, want)
		}
	})

	t.Run("SuccessWithTrailingSlach", func(t *testing.T) {
		id, err := idFromURL("https://example.com/123/")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := id, 123; got != want {
			t.Errorf("got id %d, want %d", got, want)
		}
	})

	t.Run("Error", func(t *testing.T) {
		_, err := idFromURL("https://example.com/string")

		if err == nil {
			t.Fatal("err is nil")
		}
	})
}
