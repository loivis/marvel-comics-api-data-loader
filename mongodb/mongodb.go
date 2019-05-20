package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var defaultTimeout time.Duration = 5 * time.Second

type CollectionName string

var (
	ColCharacters CollectionName = "characters"
	ColComics     CollectionName = "comics"
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

	return col.CountDocuments(ctx, bson.D{})
}

func (m *MongoDB) IncompleteCharacterIDs() ([]int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	col := m.client.Database(m.database).Collection(string(ColCharacters))

	cur, err := col.Find(ctx,
		bson.D{{Key: "intact", Value: false}},
		options.Find().SetProjection(bson.D{{Key: "id", Value: 1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("error finding incomplete characters: %v", err)
	}

	var ids []int32

	for cur.Next(ctx) {
		var elem struct{ ID int32 }
		err := cur.Decode(&elem)
		if err != nil {
			return nil, fmt.Errorf("error decoding character: %v", err)
		}

		ids = append(ids, elem.ID)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error decoding all characters: %v", err)
	}

	cur.Close(ctx)

	return ids, nil
}

func (m *MongoDB) SaveCharacter(char *m27r.Character) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	col := m.client.Database(m.database).Collection(string(ColCharacters))

	result, err := col.ReplaceOne(ctx, bson.D{{Key: "id", Value: char.ID}}, char)
	if err != nil {
		return err
	}

	log.Info().Interface("result", result).Int32("id", char.ID).Msg("character document replaced")

	return nil
}

func (m *MongoDB) SaveCharacters(chars []*m27r.Character) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	ids, err := m.getAllIds(ctx, ColCharacters)
	if err != nil {
		return err
	}

	newChars := diffCharacters(ids, chars)
	if len(newChars) == 0 {
		log.Info().Msg("no new characters to save")
		return nil
	}

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

func (m *MongoDB) SaveComics(comics []*m27r.Comic) error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	ids, err := m.getAllIds(ctx, ColComics)
	if err != nil {
		return err
	}

	newComics := diffComics(ids, comics)

	if len(newComics) == 0 {
		log.Info().Msg("no new comics to save")
		return nil
	}

	docs := []interface{}{}
	for _, comic := range newComics {
		docs = append(docs, comic)
	}

	col := m.client.Database(m.database).Collection(string(ColComics))
	_, err = col.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	log.Info().Int("count", len(docs)).Msg("new comics saved")

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

	cur, err := col.Find(ctx, bson.D{})
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

func diffCharacters(ids []int32, chars []*m27r.Character) []*m27r.Character {
	m := make(map[int32]struct{}, len(ids))
	for i := range ids {
		m[ids[i]] = struct{}{}
	}

	r := []*m27r.Character{}
	for i := range chars {
		if _, ok := m[chars[i].ID]; ok {
			continue
		}
		r = append(r, chars[i])
	}

	return r
}

func diffComics(ids []int32, comics []*m27r.Comic) []*m27r.Comic {
	log.Info().Int("count", len(ids)).Int("comics", len(comics)).Msg("diff comics")
	m := make(map[int32]struct{}, len(ids))
	for i := range ids {
		m[ids[i]] = struct{}{}
	}

	r := []*m27r.Comic{}
	for i := range comics {
		if _, ok := m[comics[i].ID]; ok {
			continue
		}
		r = append(r, comics[i])
	}

	return r
}
