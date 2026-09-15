package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const counterSlug = "requests-count"

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
		CREATE TABLE IF NOT EXISTS counters (
			id    bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			title text   NOT NULL UNIQUE,
			value bigint NOT NULL
		)`
	if _, err := s.db.ExecContext(ctx, schema); err != nil {
		return err
	}

	const seed = `
		INSERT INTO counters (title, value) VALUES ($1, 0)
		ON CONFLICT (title) DO NOTHING`
	_, err := s.db.ExecContext(ctx, seed, counterSlug)
	return err
}

func (s *Store) Increment(ctx context.Context) (int64, error) {
	const q = `
		INSERT INTO counters (title, value) VALUES ($1, 1)
		ON CONFLICT (title) DO UPDATE SET value = counters.value + 1
		RETURNING value`

	var n int64
	if err := s.db.QueryRowContext(ctx, q, counterSlug).Scan(&n); err != nil {
		return 0, fmt.Errorf("increment counter: %w", err)
	}
	return n, nil
}

func (s *Store) Count(ctx context.Context) (int64, error) {
	const q = `SELECT value FROM counters WHERE title = $1`

	var n int64
	err := s.db.QueryRowContext(ctx, q, counterSlug).Scan(&n)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return 0, nil
	case err != nil:
		return 0, fmt.Errorf("read counter: %w", err)
	}
	return n, nil
}
