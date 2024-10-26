package main

import (
	"flag"
	"strconv"

	"github.com/protomem/chatik/cmd/migrator/app"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

var (
	_db     = flag.String("db", "", "connection string to database")
	_action = flag.String("action", "up", "action to perform")
	_step   = flag.Int("step", 0, "number of migrations to run")
)

func init() {
	flag.Parse()
}

func main() {
	var args []string
	if *_step != 0 {
		args = append(args, strconv.Itoa(*_step))
	}

	app.Run(app.Features{
		Database: *_db,
		Action:   *_action,
		Args:     args,
	})
}
