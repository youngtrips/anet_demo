package db

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"protocol"
)

func AccountAuth(user string, pass string) (protocol.ERROR, int32) {
	row := _dbSession.QueryRow(
		"SELECT id, password, salt FROM account WHERE username=?", user)
	var id int32
	var originPass string
	var salt string
	err := row.Scan(&id, &originPass, &salt)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("invalid username")
		return protocol.ERROR_NO_FOUND_USER, 0
	case err != nil:
		log.Fatal(err)
	default:
		checksum := fmt.Sprintf("%x", md5.Sum([]byte(pass+salt)))
		if checksum == originPass {
			log.Printf("auth ok...")
		} else {
			log.Printf("wrong password...")
			return protocol.ERROR_INVALID_PASSWORD, 0
		}
	}
	return protocol.ERROR_SUCCESS, id
}
