package marvel

import (
	"net/http"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ae := &APIError{Code: http.StatusBadRequest, Message: "foo"}

		if got, want := ae.Error(), "status 400: foo"; got != want {
			t.Errorf("got error %q, want %q", got, want)
		}
	})
}

func TestPathError_Error(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ae := &PathError{Missing: []string{"foo", "bar"}}

		if got, want := ae.Error(), "missing in path: foo, bar"; got != want {
			t.Errorf("got error %q, want %q", got, want)
		}
	})
}
