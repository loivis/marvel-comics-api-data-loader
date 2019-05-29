package gcpfirestore

import (
	"testing"

	"github.com/loivis/marvel-comics-api-data-loader/m27r"
)

// doc implements m27r.Doc.
type doc struct {
	ID     int
	Intact bool
}

func (doc *doc) Identify() int {
	return doc.ID
}

func TestDiff(t *testing.T) {
	for _, tc := range []struct {
		desc string
		ids  []int
		docs []m27r.Doc
		out  []m27r.Doc
	}{
		{
			desc: "NoDiff",
			ids:  []int{1, 2, 3, 4},
			docs: []m27r.Doc{&doc{ID: 1}, &doc{ID: 2}, &doc{ID: 3}, &doc{ID: 4}},
			out:  []m27r.Doc{},
		},
		{
			desc: "LessIncoming",
			ids:  []int{1, 2, 3, 4},
			docs: []m27r.Doc{&doc{ID: 1}, &doc{ID: 2}},
			out:  []m27r.Doc{},
		},
		{
			desc: "MoreInt",
			ids:  []int{1, 2, 3, 4, 5, 6},
			docs: []m27r.Doc{&doc{ID: 1}, &doc{ID: 2}, &doc{ID: 8}, &doc{ID: 3}, &doc{ID: 4}, &doc{ID: 7}, &doc{ID: 9}},
			out:  []m27r.Doc{&doc{ID: 8}, &doc{ID: 7}, &doc{ID: 9}},
		},
		{
			desc: "WithDuplicates",
			ids:  []int{1, 2, 3},
			docs: []m27r.Doc{&doc{ID: 1}, &doc{ID: 2}, &doc{ID: 4}, &doc{ID: 3}, &doc{ID: 4}},
			out:  []m27r.Doc{&doc{ID: 4}},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			gotDocs := diff(tc.ids, tc.docs)

			if got, want := len(gotDocs), len(tc.out); got != want {
				t.Fatalf("[%s] got %d docs, want %d", tc.desc, got, want)
			}

			for i, doc := range tc.out {
				if got, want := gotDocs[i].Identify(), doc.Identify(); got != want {
					t.Errorf("[%s] got docs[%d] %d, want %d", tc.desc, i, got, want)
				}

			}
		})
	}
}

func BenchmarkDiff(b *testing.B) {
	ids := []int{}
	var docs []m27r.Doc
	for i := 0; i < 1000; i++ {
		ids = append(ids, i)
		docs = append(docs, &doc{ID: i})
	}

	for i := 0; i < b.N; i++ {
		diff(ids, docs)
	}
	// 100: BenchmarkDiff-4   	  200000	      6282 ns/op	    1046 B/op	       6 allocs/op
	// 1000: BenchmarkDiff-4   	   30000	     48931 ns/op	   13721 B/op	       6 allocs/op
	// 2000: BenchmarkDiff-4   	   20000	     93117 ns/op	   27531 B/op	       7 allocs/op
	// 5000: BenchmarkDiff-4   	    5000	    275893 ns/op	   58809 B/op	      10 allocs/op
	// 10000: BenchmarkDiff-4       3000	    480993 ns/op	  109334 B/op	      11 allocs/op
}
