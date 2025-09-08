package db

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"immodi/novel-site/internal/db/repositories"
	schema "immodi/novel-site/internal/db/schema"

	_ "modernc.org/sqlite"
)

// NewDBService creates a new database service instance
func NewDBService(databasePath string) (*DBService, error) {
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Limit to 1 writer connection (SQLite safe setup)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(0)

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Enable WAL mode (non-blocking reads during writes)
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode=WAL;"); err != nil {
		db.Close()
		cancel()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Optional: improve performance with normal sync
	if _, err := db.ExecContext(ctx, "PRAGMA synchronous=NORMAL;"); err != nil {
		db.Close()
		cancel()
		return nil, fmt.Errorf("failed to set synchronous: %w", err)
	}

	// create tables
	if _, err := db.ExecContext(ctx, schema.Schema); err != nil {
		db.Close()
		cancel()
		return nil, err
	}

	queries := repositories.New(db)

	return &DBService{
		db:      db,
		queries: queries,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

// Close gracefully shuts down the database service
func (s *DBService) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Cancel context first
	if s.cancel != nil {
		s.cancel()
	}

	// Close database connection
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}

// Generic helper function for type-safe results
func ExecuteWithResult[T any](s *DBService, fn func(context.Context, *repositories.Queries) (T, error)) (T, error) {
	s.mu.RLock()
	ctx := s.ctx
	queries := s.queries
	s.mu.RUnlock()

	return fn(ctx, queries)
}

func Execute(s *DBService, fn func(context.Context, *repositories.Queries) error) error {
	s.mu.RLock()
	ctx := s.ctx
	queries := s.queries
	s.mu.RUnlock()

	return fn(ctx, queries)
}

func ExecuteTx(s *DBService, fn func(context.Context, *repositories.Queries) error) error {
	s.mu.RLock()
	ctx := s.ctx
	s.mu.RUnlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := s.queries.WithTx(tx)

	if err := fn(ctx, q); err != nil {
		return err
	}

	return tx.Commit()
}

// Helper functions for creating sql.Null types
// NewNullString creates a sql.NullString
func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

// NewNullInt64 creates a sql.NullInt64
func NewNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

// GetStringValue safely gets string value from sql.NullString
func GetStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// GetInt64Value safely gets int64 value from sql.NullInt64
func GetInt64Value(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}
