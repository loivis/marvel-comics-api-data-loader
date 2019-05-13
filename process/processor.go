package process

import (
	"context"

	"github.com/loivis/mcapi-loader/mcapi"
	"github.com/loivis/mcapi-loader/mcapiloader"
)

type Processor struct {
	client *mcapi.Client
	store  mcapiloader.Store
}

func NewProcessor(c *mcapi.Client, s mcapiloader.Store) *Processor {
	return &Processor{
		client: c,
		store:  s,
	}
}

func Process(ctx context.Context) error {
	return nil
}
