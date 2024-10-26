package app

import (
	"errors"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/protomem/chatik/assets"
	"github.com/protomem/chatik/pkg/syslog"
)

func Run(fs Features) {
	syslog.Log(syslog.Info).Printf("Running migrations...")
	defer syslog.Log(syslog.Info).Printf("Migrations complete")

	if fs.Database == "" {
		syslog.Log(syslog.Error).Panicf("No database connection string provided")
	}

	iofsDriver, err := iofs.New(assets.Assets, "migrations")
	if err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to load migrations: %s", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", iofsDriver, fs.Database)
	if err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to initialize migrator: %s", err)
	}

	switch strings.ToLower(fs.Action) {
	case "up":
		if len(fs.Args) == 0 {
			if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				syslog.Log(syslog.Error).Panicf("Failed to run up migrations: %s", err)
			}
		} else {
			step, err := strconv.Atoi(fs.Args[0])
			if err != nil {
				syslog.Log(syslog.Error).Panicf("Failed to parse step: %s", err)
			}

			if err := migrator.Steps(step); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				syslog.Log(syslog.Error).Panicf("Failed to run up migrations: %s", err)
			}
		}
	case "down":
		if len(fs.Args) == 0 {
			if err := migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				syslog.Log(syslog.Error).Panicf("Failed to run down migrations: %s", err)
			}
		} else {
			step, err := strconv.Atoi(fs.Args[0])
			if err != nil {
				syslog.Log(syslog.Error).Panicf("Failed to parse step: %s", err)
			}

			if err := migrator.Steps(-step); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				syslog.Log(syslog.Error).Panicf("Failed to run down migrations: %s", err)
			}
		}
	case "drop":
		if err := migrator.Drop(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			syslog.Log(syslog.Error).Panicf("Failed to drop migrations: %s", err)
		}
	case "goto":
		if len(fs.Args) == 0 {
			syslog.Log(syslog.Error).Panicf("No step provided")
		}

		step, err := strconv.Atoi(fs.Args[0])
		if err != nil {
			syslog.Log(syslog.Error).Panicf("Failed to parse step: %s", err)
		}

		if err := migrator.Migrate(uint(step)); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			syslog.Log(syslog.Error).Panicf("Failed to run goto migrations: %s", err)
		}
	default:
		syslog.Log(syslog.Error).Panicf("Unknown action: %s", fs.Action)
	}

	if _, err := migrator.Close(); err != nil {
		syslog.Log(syslog.Error).Panicf("Failed to close migrator: %s", err)
	}
}
