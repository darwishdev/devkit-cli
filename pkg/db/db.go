package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type DbInterface interface {
	Open(dbSource string) (*sql.DB, error)
}
type Db struct {
}

func NewDb() DbInterface {
	return &Db{}
}

func (d *Db) Open(dbSource string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
