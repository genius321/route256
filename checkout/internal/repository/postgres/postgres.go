package postgres

import (
	"route256/libs/postgres/tx"

	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	provider tx.DBProvider
}

func New(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}
