package process

import (
	"context"
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/loivis/mcapi-loader/marvel/mclient"
	"github.com/loivis/mcapi-loader/mcapiloader"
)

const marvelComicsAPIPrefix = "http://gateway.marvel.com/v1/public/"

type Processor struct {
	mclient    *mclient.Marvel
	privateKey string
	publicKey  string
	timeout    time.Duration
	limit      int32

	store mcapiloader.Store

	concurrency int
}

func NewProcessor(mc *mclient.Marvel, s mcapiloader.Store, private, public string) *Processor {
	return &Processor{
		mclient:    mc,
		privateKey: private,
		publicKey:  public,
		timeout:    30 * time.Second,
		limit:      100,

		store: s,

		concurrency: 20,
	}
}

func (p *Processor) Process(ctx context.Context) error {
	err := p.loadCharacters(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) setParams(ctx context.Context, params mcapiloader.Params) {
	ts := fmt.Sprintf("%d", time.Now().Unix())
	params.SetApikey(p.publicKey)
	params.SetContext(ctx)
	params.SetHash(fmt.Sprintf("%x", md5.Sum([]byte(ts+p.privateKey+p.publicKey))))
	params.SetTs(ts)
	params.SetTimeout(p.timeout)
}

func idFromURL(in string) (int32, error) {
	ss := strings.Split(strings.Trim(in, "/"), "/")
	s := ss[len(ss)-1]
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}
