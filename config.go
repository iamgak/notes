package main

import (
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func InitRedis(name, passw string) *redis.Client {

	redis_port := 6379
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", name, redis_port),
		Password: passw, // no password set
		DB:       0,     // use default DB
	})

	return client
}
