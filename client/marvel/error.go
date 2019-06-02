package marvel

import (
	"fmt"
	"strings"
)

type APIError struct {
	Code    int
	Message string
}

func (ae *APIError) Error() string {
	return fmt.Sprintf("status %d: %s", ae.Code, ae.Message)
}

type PathError struct {
	Missing []string
}

func (pe *PathError) Error() string {
	return fmt.Sprintf("missing in path: %s", strings.Join(pe.Missing, ", "))
}
