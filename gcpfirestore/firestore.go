package gcpfirestore

import (
	"context"
	"fmt"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/loivis/marvel-comics-api-data-loader/m27r"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	ColCharacters = "characters"
	ColComics     = "comics"
	ColCreators   = "creators"
	ColEvents     = "events"
	ColSeries     = "series"
	ColStories    = "stories"
)

type Firestore struct {
	client      *firestore.Client
	concurrency int
}

func New(ctx context.Context, credentials []byte) (*Firestore, error) {
	option := option.WithCredentialsJSON(credentials)
	app, err := firebase.NewApp(ctx, nil, option)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	return &Firestore{
		client:      client,
		concurrency: 20,
	}, nil
}

func (s *Firestore) GetCount(ctx context.Context, collection string) (int, error) {
	refs, err := s.client.Collection(collection).DocumentRefs(ctx).GetAll()
	if err != nil {
		return 0, err
	}

	return len(refs), nil
}

func (s *Firestore) IncompleteIDs(ctx context.Context, collection string) ([]int, error) {
	var ids []int

	iter := s.client.Collection(collection).Where("intact", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var s struct{ ID int }
		if err := doc.DataTo(&s); err != nil {
			return nil, err
		}
		ids = append(ids, s.ID)
	}

	return ids, nil
}

func (s *Firestore) SaveCharacters(ctx context.Context, chars []*m27r.Character) error {
	var docs []m27r.Doc

	for _, doc := range chars {
		docs = append(docs, doc)
	}

	return s.saveAll(ctx, ColCharacters, docs)
}

func (s *Firestore) SaveComics(ctx context.Context, comics []*m27r.Comic) error {
	var docs []m27r.Doc

	for _, doc := range comics {
		docs = append(docs, doc)
	}

	return s.saveAll(ctx, ColComics, docs)
}

func (s *Firestore) SaveCreators(ctx context.Context, creators []*m27r.Creator) error {
	var docs []m27r.Doc

	for _, doc := range creators {
		docs = append(docs, doc)
	}

	return s.saveAll(ctx, ColCreators, docs)
}

func (s *Firestore) SaveEvents(ctx context.Context, events []*m27r.Event) error {
	var docs []m27r.Doc

	for _, doc := range events {
		docs = append(docs, doc)
	}

	return s.saveAll(ctx, ColEvents, docs)
}

func (s *Firestore) SaveSeries(ctx context.Context, series []*m27r.Series) error {
	var docs []m27r.Doc

	for _, doc := range series {
		docs = append(docs, doc)
	}

	return s.saveAll(ctx, ColSeries, docs)
}

func (s *Firestore) SaveStories(ctx context.Context, stories []*m27r.Story) error {
	var docs []m27r.Doc

	for _, doc := range stories {
		docs = append(docs, doc)
	}

	return s.saveAll(ctx, ColStories, docs)
}

func (s *Firestore) SaveOne(ctx context.Context, doc m27r.Doc) error {
	var collection string

	switch doc.(type) {
	case *m27r.Character:
		collection = ColCharacters
	case *m27r.Comic:
		collection = ColComics
	case *m27r.Creator:
		collection = ColCreators
	case *m27r.Event:
		collection = ColEvents
	case *m27r.Series:
		collection = ColSeries
	case *m27r.Story:
		collection = ColStories
	default:
		return fmt.Errorf("unsupported type: %T", doc)
	}

	return s.save(ctx, collection, doc)
}

func (s *Firestore) saveAll(ctx context.Context, collection string, docs []m27r.Doc) error {
	ids, err := s.allIDs(ctx, collection)
	if err != nil {
		return fmt.Errorf("failed to get all ids in collection %s: %v", collection, err)
	}

	newDocs := diff(ids, docs)

	if len(newDocs) == 0 {
		log.Info().Msg("no new docs to save")
		return nil
	}

	conCh := make(chan struct{}, s.concurrency)

	for _, doc := range newDocs {
		conCh <- struct{}{}

		doc := doc

		go func() {
			_ = s.save(ctx, collection, doc)
			<-conCh
		}()
	}

	for i := 0; i < s.concurrency; i++ {
		conCh <- struct{}{}
	}

	return nil
}

func (s *Firestore) save(ctx context.Context, collection string, doc m27r.Doc) error {
	id := strconv.FormatInt(int64(doc.Identify()), 10)

	_, err := s.client.Collection(collection).Doc(id).Set(ctx, doc)
	if err != nil {
		log.Error().Str("collection", collection).Str("id", id).Msg("failed to save document")
		return err
	}

	log.Info().Str("collection", collection).Str("id", id).Msg("saved document")

	return nil
}

func (s *Firestore) allIDs(ctx context.Context, collection string) ([]int, error) {
	var ids []int

	iter := s.client.Collection(collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var s struct{ ID int }
		if err := doc.DataTo(&s); err != nil {
			return nil, err
		}
		ids = append(ids, s.ID)
	}

	return ids, nil
}

func diff(ids []int, docs []m27r.Doc) []m27r.Doc {
	m := make(map[int]struct{}, len(ids))
	for i := range ids {
		m[ids[i]] = struct{}{}
	}

	var r []m27r.Doc
	seen := map[int]struct{}{}
	for i := range docs {
		if _, ok := m[docs[i].Identify()]; ok {
			continue
		}

		if _, ok := seen[docs[i].Identify()]; ok {
			continue
		}

		seen[docs[i].Identify()] = struct{}{}
		r = append(r, docs[i])
	}

	return r
}
