package m27r

import (
	"context"
	"time"
)

// Store .
type Store interface {
	GetCount(ctx context.Context, collection string) (int, error)
	IncompleteIDs(ctx context.Context, collection string) ([]int, error)
	SaveCharacters(ctx context.Context, chars []*Character) error
	SaveComics(ctx context.Context, comics []*Comic) error
	SaveCreators(ctx context.Context, creators []*Creator) error
	SaveEvents(ctx context.Context, events []*Event) error
	SaveSeries(ctx context.Context, series []*Series) error
	SaveStories(ctx context.Context, stories []*Story) error
	SaveOne(ctx context.Context, doc Doc) error
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
	Identify() int
}
