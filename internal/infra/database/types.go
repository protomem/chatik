package database

import "github.com/Masterminds/squirrel"

type QueryBuilder = squirrel.StatementBuilderType

type Scanner interface {
	Scan(dest ...any) error
}
