package db

import (
	"database/sql"
	"log"
)

type User struct {
	Id   int32
	Name string
}

func LoadUser(id int32) *User {
	row := _dbSession.QueryRow(
		"SELECT id, name FROM user WHERE id=?", id)
	user := new(User)
	err := row.Scan(&user.Id, &user.Name)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		log.Fatal(err)
	default:
		break
	}
	return user
}
