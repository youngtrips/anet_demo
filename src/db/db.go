package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	_dbSession *sql.DB
)

func init() {
	_dbSession = nil
}

func Open(host string, port int16, user string, pass string, dbname string) {
	baseDSN := ""
	if user != "" {
		baseDSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, pass, host, port, dbname)
	} else {
		baseDSN = fmt.Sprintf("tcp(%s:%d)/%s", user, pass, host, port, dbname)
	}

	var err error
	_dbSession, err = sql.Open("mysql", baseDSN)
	if err != nil {
		log.Printf("open database failed: %s", err.Error())
	}
}

func Close() {
	if _dbSession != nil {
		_dbSession.Close()
		_dbSession = nil
	}
}
