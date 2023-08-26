package main

import (
	"flag"
	"log"
	"os"

	"github.com/protomem/chatik/internal/api"
	"github.com/protomem/chatik/internal/config"
)

var _confFile string

func init() {
	flag.StringVar(&_confFile, "conf", "", "path to config file")
}

func main() {
	var err error

	flag.Parse()

	if _confFile != "" {
		err := config.Load(_confFile)
		if err != nil {
			log.Printf("error: %s", err.Error())
			os.Exit(1)
		}
	}

	conf, err := config.Parse()
	if err != nil {
		log.Printf("error: %s", err.Error())
		os.Exit(1)
	}

	srv := api.NewServer(conf)
	err = srv.Run()
	if err != nil {
		log.Printf("error: %s", err.Error())
		os.Exit(1)
	}
}
