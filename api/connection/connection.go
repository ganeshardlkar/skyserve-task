package connection

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "username"
	password = "password"
	dbname   = "db_name"
)

func ConnectToDB() *sqlx.DB {

	dbUrl := fmt.Sprintf("user=%s port=%d dbname=%s sslmode=disable password=%s host=localhost", user, port, dbname, password)

	db, err := sqlx.Connect("postgres", dbUrl)
	if err != nil {
		log.Println(err)
		log.Fatal("Failed to connect to DB")
		return nil
	}

	return db
}
