package main

import (
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

type Opts struct {
	Debug string
}

func init() {
	log.SetLevel(log.DebugLevel)
}

func main() {
	opts := &Opts{}
	args, err := flags.Parse(opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, arg := range args {
		shorten(arg)
	}

	log.Println("Done!")
}
