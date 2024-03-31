package models

import (
	"log"
	"skyserver/structs"

	"github.com/jmoiron/sqlx"
)

func CreateUser(db *sqlx.DB, user structs.User) error {
	log.Println("Creating new user")
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	return err
}

func GetUsers(db *sqlx.DB) ([]*structs.User, error) {
	var users []*structs.User
	err := db.Select(&users, "SELECT * FROM users")
	return users, err
}

func GetUserByUsername(db *sqlx.DB, username string) (structs.User, error) {
	var user structs.User
	err := db.Select(&user, "SELECT * FROM users where username=?", username)
	return user, err
}
