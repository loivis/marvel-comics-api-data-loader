package m27r

import (
	"context"
	"time"
)

type Store interface {
	GetCount(collection string) (count int64, err error)
	// SaveIDs(collection string, ids []int32) error
	IncompleteCharacterIDs() ([]int32, error)
	SaveCharacter(char *Character) error
	SaveCharacters(chars []*Character) error
}

type Params interface {
	SetApikey(string)
	SetContext(context.Context)
	SetHash(string)
	SetTs(string)
	SetTimeout(time.Duration)
}
