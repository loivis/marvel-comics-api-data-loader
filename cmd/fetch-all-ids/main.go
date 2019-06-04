package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/avast/retry-go"
	flag "github.com/spf13/pflag"

	"github.com/loivis/marvel-comics-api-data-loader/client/marvel"
	"github.com/loivis/marvel-comics-api-data-loader/maco"
)

var mTypeFunc = map[string]func(int, int) ([]int, error){
	maco.TypeCharacters: fetchCharacterIDs,
	maco.TypeComics:     fetchComicIDs,
	maco.TypeCreators:   fetchCreatorIDs,
	maco.TypeEvents:     fetchEventIDs,
	maco.TypeSeries:     fetchSeriesIDs,
	maco.TypeStories:    fetchStoryIDs,
}

// variables for commandline flags
var (
	filePath   string
	privateKey string
	publicKey  string
)

// global variables
var (
	ctx     context.Context
	err     error
	file    *os.File
	mStrCol map[string]*collection
	mc      *marvel.Client
)

const (
	baseURL = "https://gateway.marvel.com/v1/public/"
	limit   = 100
)

type collection struct {
	Offset int   `json:"offset,omitempty"`
	IDs    []int `json:"ids,omitempty"`
}

func init() {
	flag.StringVar(&filePath, "file-path", "", "path to id mapping file")
	flag.StringVar(&privateKey, "private-key", "", "private key for marvel comics api")
	flag.StringVar(&publicKey, "public-key", "", "public key for marvel comics api")
	flag.Parse()

	ctx = context.Background()
	mc = marvel.NewClient(baseURL, privateKey, publicKey)
}

func main() {
	if filePath == "" || privateKey == "" || publicKey == "" {
		fmt.Println("Please provide all flags below:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to open file %q: %v", filePath, err)
	}
	defer file.Close()

	mStrCol, err = current(file)
	if err != nil {
		log.Fatalf("error reading current mappings: %v", err)
	}

	log.Printf("current mappings: %v", mStrCol)

	for typ, f := range mTypeFunc {
		err := fetchIDs(ctx, typ, f)
		if err != nil {
			log.Fatalf("error fetching %s ids: %v", typ, err)
		}
	}

	log.Println("DONE")
}

func current(r io.Reader) (map[string]*collection, error) {
	m := make(map[string]*collection)

	err := json.NewDecoder(r).Decode(&m)

	if err != nil {
		if err.Error() != "EOF" {
			return nil, err
		}
		err = nil
	}

	for typ := range mTypeFunc {
		if _, ok := m[typ]; !ok {
			m[typ] = &collection{}
		}
	}

	return m, err
}

func fetchIDs(ctx context.Context, typ string, f func(int, int) ([]int, error)) error {
	count, err := mc.GetCount(ctx, typ)
	if err != nil {
		return err
	}

	log.Printf("total %d %s", count, typ)

	offset := mStrCol[typ].Offset

	for i := offset / limit; i < count/limit+1; i++ {
		log.Printf("fetching %s offset %d", typ, offset)
		retry.Do(
			func() error {
				ids, err := f(offset, limit)
				if err != nil {
					return err
				}

				mStrCol[typ].Offset = offset + limit
				mStrCol[typ].IDs = append(mStrCol[typ].IDs, ids...)

				return nil
			},
			retry.RetryIf(retryIf(offset)),
		)

		b, _ := json.Marshal(mStrCol)

		_, err = file.WriteAt(b, 0)
		if err != nil {
			log.Fatalf("error save mappings: %v", err)
		}

		log.Printf("saved %s offset %d", typ, offset)

		offset += limit
	}

	return nil
}

func fetchCharacterIDs(offset, limit int) ([]int, error) {
	res, err := mc.GetCharacters(ctx, &marvel.Params{Offset: offset, Limit: 100})
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(res))
	for i, elem := range res {
		ids[i] = elem.ID
	}

	return ids, nil
}

func fetchComicIDs(offset, limit int) ([]int, error) {
	res, err := mc.GetComics(ctx, &marvel.Params{Offset: offset, Limit: 100})
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(res))
	for i, elem := range res {
		ids[i] = elem.ID
	}

	return ids, nil
}

func fetchCreatorIDs(offset, limit int) ([]int, error) {
	res, err := mc.GetCreators(ctx, &marvel.Params{Offset: offset, Limit: 100})
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(res))
	for i, elem := range res {
		ids[i] = elem.ID
	}

	return ids, nil
}

func fetchEventIDs(offset, limit int) ([]int, error) {
	res, err := mc.GetEvents(ctx, &marvel.Params{Offset: offset, Limit: 100})
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(res))
	for i, elem := range res {
		ids[i] = elem.ID
	}

	return ids, nil
}

func fetchSeriesIDs(offset, limit int) ([]int, error) {
	res, err := mc.GetSeries(ctx, &marvel.Params{Offset: offset, Limit: 100})
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(res))
	for i, elem := range res {
		ids[i] = elem.ID
	}

	return ids, nil
}

func fetchStoryIDs(offset, limit int) ([]int, error) {
	res, err := mc.GetStories(ctx, &marvel.Params{Offset: offset, Limit: 100})
	if err != nil {
		return nil, err
	}

	ids := make([]int, len(res))
	for i, elem := range res {
		ids[i] = elem.ID
	}

	return ids, nil
}

func retryIf(offset int) func(error) bool {
	return func(err error) bool {
		if v, ok := err.(*marvel.APIError); ok && v.Code != 429 {
			log.Printf("[offset %d] retryable %T error: status code %d", offset, err, v.Code)
			return true
		}

		if err, ok := err.(net.Error); ok && err.Timeout() {
			log.Printf("[offset %d] retryable %[2]T error: %[2]v", offset, err)
			return true
		}

		log.Printf("[offset %d] non-retryable %[2]T error: %[2]v", offset, err)

		return false
	}
}
