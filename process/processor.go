package process

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/loivis/marvel-comics-api-data-loader/client/marvel"
	"github.com/loivis/marvel-comics-api-data-loader/maco"
)

type Processor struct {
	mclient    *marvel.Client
	privateKey string
	publicKey  string
	timeout    time.Duration
	limit      int

	store      maco.Store
	storeBatch int

	concurrency int
}

func NewProcessor(mc *marvel.Client, s maco.Store, private, public string) *Processor {
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

	err = p.loadCharacters(ctx)
	if err != nil {
		return err
	}

	err = p.loadComics(ctx)
	if err != nil {
		return err
	}

	err = p.loadCreators(ctx)
	if err != nil {
		return err
	}

	err = p.loadEvents(ctx)
	if err != nil {
		return err
	}

	err = p.loadSeries(ctx)
	if err != nil {
		return err
	}

	err = p.loadStories(ctx)
	if err != nil {
		return err
	}

	return nil
}

func idFromURL(in string) (int, error) {
	ss := strings.Split(strings.Trim(in, "/"), "/")
	s := ss[len(ss)-1]
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func retryIf(offset int) func(error) bool {
	return func(err error) bool {
		if v, ok := err.(*marvel.APIError); ok && v.Code != 429 {
			log.Error().Int("offset", offset).Int("code", v.Code).Msgf("retryable api error: %v", err)
			return true
		}

		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Error().Int("offset", offset).Msgf("retryable timeout %[1]T error: %[1]v", err)
			return true
		}

		log.Error().Int("offset", offset).Msgf("non-retryable %[1]T error: %[1]v", err)

		return false
	}
}

func retryLog(offset int) func(uint, error) {
	return func(n uint, err error) {
		log.Info().Int("offset", offset).Uint("n", n).Msgf("retry on %[1]T error: %[1]v", err)
	}
}
