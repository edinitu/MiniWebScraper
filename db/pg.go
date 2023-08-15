package db

import (
	"database/sql"
	"go_prj/product"
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

type DbConn struct {
	Db *sql.DB
}

func InitDb(config PgConfig) (*DbConn, error) {
	var connStr = "user=" + config.User + " password=" + config.Password +
		" dbname=" + config.Dbname + " sslmode=" + config.Sslmode
	db, err := sql.Open(POSTGRE_DRIVER, connStr)
	if err != nil {
		log.Println("Could not open database connection")
		return nil, err
	}
	return &DbConn{Db: db}, nil
}

func InsertProducts(products map[string]product.Product, conn *DbConn) {
	st, err := conn.Db.Prepare("INSERT INTO nssx.services(service_name) values ($1) RETURNING service_id")
	if err != nil {
		log.Println(err)
	}
	for _, val := range products {
		var id int = 0
		err := st.QueryRow(val.Name).Scan(&id)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		log.Printf("Insert done, id:%d", id)
	}
}
