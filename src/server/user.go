package server

import (
	"anet"
	"db"
)

type UserSession struct {
	Id      int32
	User    *db.User
	Session *anet.Session
}
