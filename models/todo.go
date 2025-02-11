package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type ToDo struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Visibility bool      `json:"visibility"`
	Editable   bool      `json:"editable"`
	Deleted    bool      `json:"deleted"`
	Updated    bool      `json:"updated"`
	Version    int       `json:"version"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ToDoModel struct {
	db    *sql.DB
	redis *redis.Client
}

func (m *ToDoModel) Close() {
	m.redis.Close()
	m.db.Close()
}

func NewModels(db *sql.DB, redis *redis.Client) *ToDoModel {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return &ToDoModel{
		db:    db,
		redis: redis,
		// cancel: cancel,
		// ctx:    ctx,
	}
}

func (c *ToDoModel) CreateTodo(ctx context.Context, todo *ToDo) error {
	query := `INSERT INTO todo (title, content, visibility, editable, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)`

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, todo.Title, todo.Content, todo.Visibility, todo.Editable, time.Now(), time.Now())
	return err
}

func (c *ToDoModel) UpdateTodo(ctx context.Context, id int, todo *ToDo) error {
	query := `UPDATE todo
			  SET title = ?, content = ?, visibility = ?, editable = ?, updated_at = ?
			  WHERE id = ?`

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, todo.Title, todo.Content, todo.Visibility, todo.Editable, time.Now(), id)
	return err
}

func (c *ToDoModel) ToDoListing(ctx context.Context) ([]*ToDo, error) {
	query := `SELECT id, title, content, visibility, editable, deleted, updated, version, created_at, updated_at 
				FROM todo`

	queryBytes, err := json.Marshal(query)
	if err != nil {
		panic(err)
	}

	listing := []*ToDo{}
	val, err := c.getRedis(ctx, string(queryBytes))
	if err == nil {
		return val, err
	}

	if err != redis.Nil {
		return nil, err
	}

	// var listing []*ToDo

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		todo := &ToDo{}
		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Content,
			&todo.Visibility,
			&todo.Editable,
			&todo.Deleted,
			&todo.Updated,
			&todo.Version,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		listing = append(listing, todo)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	err = c.setRedis(ctx, string(queryBytes), listing, 5*time.Minute)
	return listing, err
}

func (c *ToDoModel) getRedis(ctx context.Context, key string) ([]*ToDo, error) {
	var listing []*ToDo
	val, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		return listing, err
	}

	// Deserialize the cached result
	err = json.Unmarshal([]byte(val), &listing)
	return listing, err
}

func (c *ToDoModel) setRedis(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.redis.Set(ctx, key, data, expiration).Err()
	return err
}
