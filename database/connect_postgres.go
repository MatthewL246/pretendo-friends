package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/PretendoNetwork/friends-secure/globals"
)

var postgres *sql.DB

func connectPostgres() {
	var err error

	postgres, err = sql.Open("postgres", os.Getenv("DATABASE_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
	}

	initPostgresWiiU()
	// TODO: 3DS database
}