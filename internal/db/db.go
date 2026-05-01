package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/goutils"
	"github.com/joseph0x45/sad"
	"github.com/joseph0x45/tessera/internal/buildinfo"
)

type Conn struct {
	db      *sqlx.DB
	verbose bool
}

func (c *Conn) Close() {
	if c.verbose {
		log.Println("Database connection closed")
	}
	c.db.Close()
}

func GetConn(verbose bool) *Conn {
	dbPath := goutils.Setup()
	if buildinfo.Version == "debug" {
		dbPath = "db.sqlite"
	}
	db, err := sad.OpenDBConnection(sad.DBConnectionOptions{
		Reset:             false,
		EnableForeignKeys: true,
		DatabasePath:      dbPath,
	}, migrations)
	if err != nil {
		panic(err)
	}
	if verbose {
		log.Println("Connected to database file at", dbPath)
	}
	return &Conn{db: db, verbose: verbose}
}

