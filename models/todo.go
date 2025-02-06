package models

import (
	"context"
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
)

type ToDo struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Visibility  bool      `json:"visibility"`
	Editable    bool      `json:"editable"`
	Deleted     bool      `json:"deleted"`
	Updated     bool      `json:"updated"`
	Version     int       `json:"version"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ToDoDB struct {
	db *sql.DB
	// ctx    context.Context
	redis *redis.Client
	// cancel context.CancelFunc
}

func (m *ToDoDB) Close() {
	// m.cancel()
	m.redis.Close()
	m.db.Close()
}

func NewModels(db *sql.DB, redis *redis.Client) *ToDoDB {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return &ToDoDB{
		db:    db,
		redis: redis,
		// cancel: cancel,
		// ctx:    ctx,
	}
}

func (c *ToDoDB) CreateTodo(ctx context.Context, todo *ToDo) error {
	query := `INSERT INTO todo (title, description, visibility, editable, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?)`

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, todo.Title, todo.Description, todo.Visibility, todo.Editable, time.Now(), time.Now())
	// defer ctx.cancle
	return err
}

func (c *ToDoDB) UpdateTodo(ctx context.Context, id int, todo *ToDo) error {
	query := `UPDATE todo
			  SET title = ?, description = ?, visibility = ?, editable = ?, updated_at = ?
			  WHERE id = ?`

	stmt, err := c.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, todo.Title, todo.Description, todo.Visibility, todo.Editable, time.Now(), id)
	// defer c.cancel()
	return err
}

func (c *ToDoDB) ToDoListing(ctx context.Context) ([]*ToDo, error) {
	query := `SELECT id, title, content, visibility, editable, deleted, updated, version, created_at, updated_at 
				FROM todo`

	// if userID != 0 {
	// 	query += " AND user_id = ?"
	// }

	var listing []*ToDo

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
			&todo.Description,
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
	// defer c.Close()
	return listing, err
}
