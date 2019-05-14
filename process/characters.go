package process

import (
	"context"
	"log"

	"github.com/loivis/mcapi-loader/marvel/mclient/operations"
)

func (p *Processor) characters(ctx context.Context) error {
	ids, err := p.getAllCharacters(ctx)
	if err != nil {
		return err
	}

	log.Println(ids)

	return nil
}

func (p *Processor) getAllCharacters(ctx context.Context) ([]int32, error) {
	ids := []int32{}

	params := &operations.GetCharactersCollectionParams{}
	p.setParams(params)

	chars, err := p.mclient.Operations.GetCharactersCollection(params)
	if err != nil {
		return nil, err
	}

	for _, r := range chars.Payload.Data.Results {
		ids = append(ids, r.ID)
	}
	return ids, nil
}
