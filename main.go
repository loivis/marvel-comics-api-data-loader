package main

import (
	"context"
	"log"

	"github.com/loivis/mcapi-loader/mcapi"
	"github.com/loivis/mcapi-loader/process"
	"github.com/loivis/pavium-api/store"
)

func main() {
	ctx := context.Background()
	client := mcapi.NewClient(conf.mcapiBaseURL, conf.mcapiPrivateKey, conf.mcapiPublicKey)
	store := store.New()
	p := process.NewProcessor(client, store)

	if err := p.Process(ctx); err != nil {
		log.Fatal(err)
	}
}
