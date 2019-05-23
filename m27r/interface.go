package m27r

import (
	"context"
	"time"
)

type Store interface {
	GetCount(collection string) (count int64, err error)
	// SaveIDs(collection string, ids []int32) error
	IncompleteCharacterIDs() ([]int32, error)
	IncompleteIDs(collection string) ([]int32, error)
	SaveCharacter(char *Character) error
	SaveCharacters(chars []*Character) error
	SaveComic(comic *Comic) error
	SaveComics(comics []*Comic) error
	SaveCreator(creator *Creator) error
	SaveCreators(creators []*Creator) error
}

type Params interface {
	SetApikey(string)
	SetContext(context.Context)
	SetHash(string)
	SetTs(string)
	SetTimeout(time.Duration)
}
