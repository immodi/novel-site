package db

import (
	"context"
	"database/sql"
	"immodi/novel-site/internal/db/repositories"
	"sync"
)

// DBService wraps the database connection and queries
type DBService struct {
	db      *sql.DB
	ctx     context.Context
	queries *repositories.Queries
	cancel  context.CancelFunc
	mu      sync.RWMutex
}
