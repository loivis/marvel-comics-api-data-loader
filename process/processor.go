package process

import (
	"context"

	mclient "github.com/loivis/mcapi-loader/marvel/client"
	"github.com/loivis/mcapi-loader/mcapiloader"
)

type Processor struct {
	mclient *mclient.Marvel
	store   mcapiloader.Store
}

func NewProcessor(mc *mclient.Marvel, s mcapiloader.Store) *Processor {
	return &Processor{
		mclient: mc,
		store:   s,
	}
}

func (p *Processor) Process(ctx context.Context) error {
	return nil
}
