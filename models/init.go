package models

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Init struct {
	Todo  ToDoModel
	Users UserModel
	Redis RedisStruct
	// Review ReviewModel
}

func Constructor(db *sql.DB, redis *redis.Client) *Init {
	RedisClient := RedisStruct{client: redis}
	return &Init{
		Todo:  ToDoModel{db: db, redis: RedisClient},
		Users: UserModel{db: db, redis: redis},
		// Redis: RedisStruct{client: redis},
		// Review: ReviewModel{db: db, redis: rd},
	}
}
