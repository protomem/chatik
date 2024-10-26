package main

import (
	"flag"

	"github.com/protomem/chatik/cmd/api-server/app"
)

var _confFile = flag.String("config", "configs/api-server.yaml", "path to config file")

func init() {
	flag.Parse()
}

func main() {
	app.Run(app.Features{
		ConfigFile: *_confFile,
	})
}
