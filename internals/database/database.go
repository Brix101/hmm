package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func CreateConnectionPool() *sqlx.DB {
	// open and connect at the same time, panicing on error
	db := sqlx.MustConnect("sqlite3", "storage.db")
	err := db.Ping()
	if err != nil {
		panic(err)
	} else {
		log.Println("✅✅✅ Database connection established")
	}

	return db
}
