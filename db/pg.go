package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const POSTGRE_DRIVER string = "postgres"

type PgConfig struct {
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

func InitDb(config PgConfig) (*sql.DB, error) {
	var connStr = "user=" + config.User + " password=" + config.Password +
		" dbname=" + config.Dbname + " sslmode=" + config.Sslmode
	db, err := sql.Open(POSTGRE_DRIVER, connStr)
	if err != nil {
		log.Println("Could not open database connection")
		return nil, err
	}
	return db, nil
}
