package m27r

import (
	"context"
	"time"
)

// Store .
type Store interface {
	GetCount(collection string) (count int64, err error)
	IncompleteIDs(collection string) ([]int32, error)
	SaveCharacters(chars []*Character) error
	SaveComics(comics []*Comic) error
	SaveCreators(creators []*Creator) error
	SaveEvents(creators []*Event) error
	SaveSeries(creators []*Series) error
	SaveOne(doc Doc) error
}

// Params abstracts common features of all params
type Params interface {
	SetApikey(string)
	SetContext(context.Context)
	SetHash(string)
	SetTs(string)
	SetTimeout(time.Duration)
}

type Doc interface {
	Identify() int32
}
