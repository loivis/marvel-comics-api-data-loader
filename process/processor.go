package process

import (
	"context"
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
	"github.com/loivis/marvel-comics-api-data-loader/marvel/mclient"
)

type Processor struct {
	mclient    *mclient.Marvel
	privateKey string
	publicKey  string
	timeout    time.Duration
	limit      int32

	store      m27r.Store
	storeBatch int

	concurrency int
}

func NewProcessor(mc *mclient.Marvel, s m27r.Store, private, public string) *Processor {
	return &Processor{
		mclient:    mc,
		privateKey: private,
		publicKey:  public,
		timeout:    30 * time.Second,
		limit:      100,

		store:      s,
		storeBatch: 1000,

		concurrency: 10,
	}
}

func (p *Processor) Process(ctx context.Context) error {
	var err error

	// err = p.loadCharacters(ctx)
	// if err != nil {
	// 	return err
	// }

	// err = p.loadComics(ctx)
	// if err != nil {
	// 	return err
	// }

	// err = p.loadCreators(ctx)
	// if err != nil {
	// 	return err
	// }

	// err = p.loadEvents(ctx)
	// if err != nil {
	// 	return err
	// }

	err = p.loadSeries(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) setParams(ctx context.Context, params m27r.Params) {
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
