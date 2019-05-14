package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/loivis/mcapi-loader/mcapiloader"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

var defaultTimeout time.Duration = 5 * time.Second

type CollectionName string

var (
	ColCharacters CollectionName = "characters"
)

type MongoDB struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func New(uri string, database string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %v", err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %v", err)
	}

	return &MongoDB{
		client:   client,
		database: database,
		timeout:  defaultTimeout,
	}, nil
}

func (m *MongoDB) GetCount(collection string) (count int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	col := m.client.Database(m.database).Collection(collection)

	return col.CountDocuments(ctx, bsonx.Doc{})
}

func (m *MongoDB) SaveCharacters(chars []*mcapiloader.Character) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	ids, err := m.getAllIds(ctx, ColCharacters)
	if err != nil {
		return err
	}

	newChars := diff(ids, chars)

	docs := []interface{}{}
	for _, char := range newChars {
		docs = append(docs, char)
	}

	col := m.client.Database(m.database).Collection(string(ColCharacters))
	_, err = col.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	log.Info().Int("count", len(docs)).Msg("new characters saved")

	return nil
}

// func (m *MongoDB) SaveIDs(collection string, ids []int32) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
// 	defer cancel()

// 	existing, err := m.getAllIds(ctx, ColCharacters)
// 	if err != nil {
// 		return err
// 	}

// 	newIDs := diff(existing, ids)

// 	docs := []interface{}{}
// 	for _, id := range newIDs {
// 		docs = append(docs, struct{ ID int32 }{id})
// 	}

// 	col := m.client.Database(m.database).Collection(collection)
// 	_, err = col.InsertMany(ctx, docs)
// 	if err != nil {
// 		return err
// 	}

// 	log.Info().Int("count", len(docs)).Msg("new ids saved")

// 	return nil
// }

func (m *MongoDB) getAllIds(ctx context.Context, collection CollectionName) ([]int32, error) {
	ids := []int32{}

	col := m.client.Database(m.database).Collection(string(collection))

	cur, err := col.Find(ctx, bsonx.Doc{}, options.Find())
	if err != nil {
		return nil, fmt.Errorf("error finding documents: %v", err)
	}

	for cur.Next(ctx) {
		var elem struct{ ID int32 }

		err := cur.Decode(&elem)
		if err != nil {
			return nil, fmt.Errorf("error decoding document: %v", err)
		}
		ids = append(ids, elem.ID)
	}

	if err := cur.Err(); err != nil {

		return nil, fmt.Errorf("error from cursor: %v", err)
	}

	cur.Close(ctx)

	return ids, nil
}

func diff(ids []int32, chars []*mcapiloader.Character) []*mcapiloader.Character {
	m := make(map[int32]struct{}, len(ids))
	for i := range ids {
		m[ids[i]] = struct{}{}
	}

	r := []*mcapiloader.Character{}
	for i := range chars {
		if _, ok := m[chars[i].ID]; ok {
			continue
		}
		r = append(r, chars[i])
	}

	return r
}
