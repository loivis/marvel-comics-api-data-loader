package main

import (
	"context"
	"log"

	mclient "github.com/loivis/mcapi-loader/marvel/client"
	"github.com/loivis/mcapi-loader/process"
	"github.com/loivis/mcapi-loader/store"
)

func main() {
	ctx := context.Background()

	mclient := mclient.New(nil, nil)
	store := store.New("")
	p := process.NewProcessor(mclient, store)

	if err := p.Process(ctx); err != nil {
		log.Fatal(err)
	}
}
