package process

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/loivis/mcapi-loader/marvel/mclient"
	"github.com/loivis/mcapi-loader/mcapiloader"
)

type Processor struct {
	mclient    *mclient.Marvel
	store      mcapiloader.Store
	privateKey string
	publicKey  string
	timeout    time.Duration
}

func NewProcessor(mc *mclient.Marvel, s mcapiloader.Store, private, public string) *Processor {
	return &Processor{
		mclient:    mc,
		store:      s,
		privateKey: private,
		publicKey:  public,
		timeout:    10 * time.Second,
	}
}

func (p *Processor) Process(ctx context.Context) error {
	err := p.characters(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) setParams(params mcapiloader.Params) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	params.SetApikey(p.publicKey)
	params.SetHash(fmt.Sprintf("%x", md5.Sum([]byte(ts+p.privateKey+p.publicKey))))
	params.SetTs(ts)
	params.SetTimeout(p.timeout)
}
