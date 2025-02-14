package models

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Init struct {
	Todo  ToDoModel
	Users UserModel
	// Review ReviewModel
}

func Constructor(db *sql.DB, redis *redis.Client) *Init {
	return &Init{
		Todo:  ToDoModel{db: db, redis: redis},
		Users: UserModel{db: db, redis: redis},
		Gpt:   Prompt(),
		// Review: ReviewModel{db: db, redis: rd},
	}
}
