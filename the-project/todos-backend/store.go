package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Todo struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

type Store struct {
	db *sql.DB
}

func NewStore(ctx context.Context, dsn string) (*Store, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	s := &Store{db: db}
	if err := s.migrate(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate(ctx context.Context) error {
	const schema = `
		CREATE TABLE IF NOT EXISTS todos (
			id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			content    text        NOT NULL,
			done       boolean     NOT NULL DEFAULT false,
			created_at timestamptz NOT NULL DEFAULT now()
		)`
	if _, err := s.db.ExecContext(ctx, schema); err != nil {
		return err
	}

	return s.seed(ctx)
}

// seed inserts the example todos only when the table is empty.
func (s *Store) seed(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var empty bool
	if err := tx.QueryRowContext(ctx, `SELECT NOT EXISTS (SELECT 1 FROM todos)`).Scan(&empty); err != nil {
		return err
	}
	if !empty {
		return nil
	}

	seed := []string{"Learn Kubernetes", "Deploy the todos app", "Write a README"}

	const q = `INSERT INTO todos (content) VALUES ($1)`
	for _, content := range seed {
		if _, err := tx.ExecContext(ctx, q, content); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) List(ctx context.Context) ([]Todo, error) {
	const q = `
		SELECT id, content, done, created_at
		FROM todos
		ORDER BY created_at DESC, id DESC`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}
	defer rows.Close()

	todos := []Todo{}
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.Content, &t.Done, &t.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan todo: %w", err)
		}
		todos = append(todos, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}

	return todos, nil
}

func (s *Store) Add(ctx context.Context, content string) (Todo, error) {
	const q = `
		INSERT INTO todos (content)
		VALUES ($1)
		RETURNING id, content, done, created_at`

	var t Todo
	if err := s.db.QueryRowContext(ctx, q, content).Scan(&t.ID, &t.Content, &t.Done, &t.CreatedAt); err != nil {
		return Todo{}, fmt.Errorf("create todo: %w", err)
	}

	return t, nil
}
