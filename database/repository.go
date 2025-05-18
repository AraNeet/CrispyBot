package database

import (
	"fmt"

	"github.com/surrealdb/surrealdb.go"
	"github.com/surrealdb/surrealdb.go/pkg/models"
)

func CreateUser(db *surrealdb.DB, userID string) (User, error) {
	if userID == "" {
		return User{}, fmt.Errorf("discord ID isn't passed.")
	}

	user, err := surrealdb.Create[User](db, models.Table("users"), map[interface{}]interface{}{
		"DiscordID": userID,
	})
	if err != nil {
		return User{}, fmt.Errorf("Failed to create user")
	}

	return *user, nil
}
