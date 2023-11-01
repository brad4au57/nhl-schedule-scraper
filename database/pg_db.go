package database

import (
	"database/sql"
)

var DB *sql.DB

func InitPostgresDb(dataSourceName string) error {
	var err error
	// Pass to Open
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}

	return DB.Ping()
}
