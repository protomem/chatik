package main

import (
	"flag"
	"log"
	"os"

	"github.com/protomem/chatik/internal/app"
)

var _confFile string

func init() {
	flag.StringVar(&_confFile, "conf", "", "path to config file")
}

func main() {
	flag.Parse()

	if _confFile == "" {
		log.Printf("error: no config file")
		os.Exit(1)
	}

	conf, err := app.NewConfig(_confFile)
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}

	app, err := app.New(conf)
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		log.Printf("error: %v", err)
		os.Exit(1)
	}
}
