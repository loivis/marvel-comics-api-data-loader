package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/loivis/mcapi-loader/marvel/mclient"
	"github.com/loivis/mcapi-loader/process"
	"github.com/loivis/mcapi-loader/store"
)

func main() {
	conf := readConfig()
	fmt.Fprintln(os.Stderr, conf)

	ctx := context.Background()

	mc := mclient.Default

	store := store.New("")

	p := process.NewProcessor(mc, store, conf.privateKey, conf.publicKey)

	if err := p.Process(ctx); err != nil {
		log.Fatal(err)
	}
}

type config struct {
	privateKey string
	publicKey  string
}

func readConfig() *config {
	return &config{
		privateKey: os.Getenv("MARVEL_API_PRIVATE_KEY"),
		publicKey:  os.Getenv("MARVEL_API_PUBLIC_KEY"),
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
		{"MARVEL_API_PRIVATE_KEY", hideIfSet(c.privateKey)},
		{"MARVEL_API_PUBLIC_KEY", c.publicKey},
	} {
		fmt.Fprintf(w, "%s\t%v\n", e.k, e.v)
	}
	w.Flush()
	return buf.String()
}
