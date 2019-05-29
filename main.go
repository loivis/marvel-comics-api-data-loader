package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/rs/zerolog/log"

	"github.com/loivis/marvel-comics-api-data-loader/marvel/mclient"
	"github.com/loivis/marvel-comics-api-data-loader/mongodb"
	"github.com/loivis/marvel-comics-api-data-loader/process"
)

func main() {
	conf := readConfig()
	fmt.Fprintln(os.Stderr, conf)

	ctx := context.Background()

	marvelClient := mclient.Default

	mongodb, err := mongodb.New(conf.mongodbURI, conf.mongodbDatabase)
	if err != nil {
		log.Fatal().Msgf("failed to setup mongodb: %v", err)
	}

	p := process.NewProcessor(marvelClient, mongodb, conf.privateKey, conf.publicKey)

	if err := p.Process(ctx); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

type config struct {
	mongodbURI      string
	mongodbDatabase string
	privateKey      string
	publicKey       string
}

func readConfig() *config {
	return &config{
		mongodbURI:      os.Getenv("MONGODB_URI"),
		mongodbDatabase: os.Getenv("MONGODB_DATABASE"),
		privateKey:      os.Getenv("MARVEL_API_PRIVATE_KEY"),
		publicKey:       os.Getenv("MARVEL_API_PUBLIC_KEY"),
	}
}

func (c *config) String() string {
	hideIfSet := func(v interface{}) string {
		s := ""

		switch typedV := v.(type) {
		case string:
			s = typedV
		case []string:
			s = strings.Join(typedV, ",")
		case fmt.Stringer:
			if typedV != nil {
				s = typedV.String()
			}
		}

		if s != "" {
			return "<hidden>"
		}
		return ""
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 1, 4, ' ', 0)
	for _, e := range []struct {
		k string
		v interface{}
	}{
		{"MONGODB_URI", hideIfSet(c.mongodbURI)},
		{"MONGODB_DATABASE", c.mongodbDatabase},
		{"MARVEL_API_PRIVATE_KEY", hideIfSet(c.privateKey)},
		{"MARVEL_API_PUBLIC_KEY", c.publicKey},
	} {
		fmt.Fprintf(w, "%s\t%v\n", e.k, e.v)
	}
	w.Flush()
	return buf.String()
}
