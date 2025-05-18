package database

import (
	"os"

	"github.com/surrealdb/surrealdb.go"
)

var (
	host          string = os.Getenv("DBHOST")
	auth_user     string = os.Getenv("DBUSER")
	auth_password string = os.Getenv("DBPASSWORD")
	db_name       string = os.Getenv("DB")
	ns            string = os.Getenv("NAMESPACE")
)

func Connect() (session *surrealdb.DB) {
	session, err := surrealdb.New(host)
	if err != nil {
		panic(err)
	}

	if err = session.Use(ns, db_name); err != nil {
		panic(err)
	}

	authData := &surrealdb.Auth{
		Username: auth_user,
		Password: auth_password,
	}
	token, err := session.SignIn(authData)
	if err != nil {
		panic(err)
	}

	if err := session.Authenticate(token); err != nil {
		panic(err)
	}

	return session
}
