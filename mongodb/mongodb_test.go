package mongodb

import (
	"context"
	"testing"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
)

func TestMongoDB_GetCount(t *testing.T) {
	m, err := New("mongodb://localhost:27017", "marvel_test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defer m.client.Database("marvel_test").Drop(context.Background())

	type doc struct{ name string }
	docs := []interface{}{
		doc{name: "foo"},
		doc{name: "bar"},
		doc{name: "baz"},
	}

	m.client.Database("marvel_test").Collection("foo").InsertMany(
		context.Background(),
		docs,
	)
	count, err := m.GetCount("foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := count, int64(len(docs)); got != want {
		t.Errorf("got  %d, want %d", got, want)
	}
}

// func TestMongoDB_SaveIDs(t *testing.T) {
// 	m, docs, err := setupDatabase("marvel_test", "foo")
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	defer m.client.Database("marvel_test").Drop(context.Background())

// 	ids := []int32{100, 101, 102}
// 	err = m.SaveIDs("foo", ids)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	gotCount, _ := m.GetCount("foo")
// 	if got, want := int(gotCount), len(docs)+len(ids); got != want {
// 		t.Errorf("got %v docs, want %v", got, want)
// 	}
// }

func TestMongoDB_GetAllIDs(t *testing.T) {
	m, docs, err := setupDatabase("marvel_test", "foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defer m.client.Database("marvel_test").Drop(context.Background())

	si, err := m.getAllIds(context.Background(), "foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := len(si), len(docs); got != want {
		t.Errorf("got %d ids, want %d", got, want)
	}
}

func TestMongoDB_IncompleteCharacterIDs(t *testing.T) {
	m, _, err := setupDatabase("marvel_test", string(ColCharacters))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defer m.client.Database("marvel_test").Drop(context.Background())

	si, err := m.IncompleteCharacterIDs()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := len(si), 1; got != want {
		t.Errorf("got %d ids, want %d", got, want)
	}
}

func setupDatabase(database, collection string) (*MongoDB, []interface{}, error) {
	m, err := New("mongodb://localhost:27017", database)
	if err != nil {
		return nil, nil, err
	}

	if err := m.client.Database(database).Drop(context.Background()); err != nil {
		return nil, nil, err
	}

	docs := []interface{}{}
	for i := 0; i < 77; i++ {
		if i == 0 {
			docs = append(docs, &m27r.Character{
				ID:     int32(i),
				Intact: false,
			})
			continue
		}

		docs = append(docs, &m27r.Character{
			ID:     int32(i),
			Intact: true,
		})
	}

	m.client.Database(database).Collection(collection).InsertMany(
		context.Background(),
		docs,
	)

	return m, docs, nil
}

func TestDiff(t *testing.T) {
	for _, tc := range []struct {
		desc  string
		ids   []int32
		chars []*m27r.Character
		out   []*m27r.Character
	}{
		{
			desc:  "NoDiff",
			ids:   []int32{1, 2, 3, 4},
			chars: []*m27r.Character{{ID: 1}, {ID: 2}, {ID: 3}, {ID: 4}},
			out:   []*m27r.Character{},
		},
		{
			desc:  "LessIncoming",
			ids:   []int32{1, 2, 3, 4},
			chars: []*m27r.Character{{ID: 1}, {ID: 2}},
			out:   []*m27r.Character{},
		},
		{
			desc:  "MoreInt32",
			ids:   []int32{1, 2, 3, 4, 5, 6},
			chars: []*m27r.Character{{ID: 1}, {ID: 2}, {ID: 8}, {ID: 3}, {ID: 4}, {ID: 7}, {ID: 9}},
			out:   []*m27r.Character{{ID: 8}, {ID: 7}, {ID: 9}},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			gotChars := diff(tc.ids, tc.chars)

			if got, want := len(gotChars), len(tc.out); got != want {
				t.Fatalf("[%s] got %d chars, want %d", tc.desc, got, want)
			}

			for i, char := range tc.out {
				if got, want := gotChars[i].ID, char.ID; got != want {
					t.Errorf("[%s] got chars[%d] %d, want %d", tc.desc, i, got, want)
				}

			}
		})
	}
}

func BenchmarkDiff(b *testing.B) {
	s1 := []int32{}
	s2 := []*m27r.Character{}
	for i := 0; i < 5000; i++ {
		s1 = append(s1, int32(i))
		s2 = append(s2, &m27r.Character{ID: int32(i)})
	}

	for i := 0; i < b.N; i++ {
		diff(s1, s2)
	}
	// 100: BenchmarkDiff-4   	  200000	      6282 ns/op	    1046 B/op	       6 allocs/op
	// 1000: BenchmarkDiff-4   	   30000	     48931 ns/op	   13721 B/op	       6 allocs/op
	// 2000: BenchmarkDiff-4   	   20000	     93117 ns/op	   27531 B/op	       7 allocs/op
	// 5000: BenchmarkDiff-4   	    5000	    275893 ns/op	   58809 B/op	      10 allocs/op
	// 10000: BenchmarkDiff-4       3000	    480993 ns/op	  109334 B/op	      11 allocs/op
}
