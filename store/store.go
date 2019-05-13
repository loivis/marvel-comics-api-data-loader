package store

type Store struct {
	uri string
}

func New(uri string) *Store {
	return &Store{
		uri: uri,
	}
}
