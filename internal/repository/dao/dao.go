package dao

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
)

type database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *database {
	return &database{
		db: db,
	}
}
