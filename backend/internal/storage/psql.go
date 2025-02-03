package storage

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nessai1/aiinterview/internal/utils"

	"context"
	"database/sql"
	"fmt"
	"github.com/nessai1/aiinterview/internal/domain"
	"github.com/pressly/goose"
)

type PSQLStorage struct {
	db *sql.DB
}

func NewPSQLStorageFromAddr(addr string) (*PSQLStorage, error) {
	conn, err := sql.Open("pgx", addr)
	if err != nil {
		return nil, fmt.Errorf("cannot create DB connection: %w", err)
	}

	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping to created DB connection: %w", err)
	}

	s, err := NewPSQLStorage(conn)
	if err != nil {
		return nil, fmt.Errorf("cannot create PSQL storage with created conn: %w", err)
	}

	return s, nil
}

func NewPSQLStorage(db *sql.DB) (*PSQLStorage, error) {
	err := goose.Up(db, "migrations/psql")
	if err != nil {
		return nil, fmt.Errorf("cannot make migrations for PSQL storage: %w", err)
	}

	s := PSQLStorage{
		db: db,
	}

	return &s, nil
}

func (s *PSQLStorage) GetUserInterviewList(_ context.Context, _ string) ([]domain.Interview, error) {
	i := make([]domain.Interview, 0)

	return i, nil
}

func (s *PSQLStorage) RegisterUser(ctx context.Context) (domain.User, error) {
	userUUID, err := utils.GenerateUUIDv7()
	if err != nil {
		return domain.User{}, fmt.Errorf("cannot generate UUIDv7: %w", err)
	}

	_, err = s.db.ExecContext(ctx, "INSERT INTO users (uuid) VALUES ($1)", userUUID)

	if err != nil {
		return domain.User{}, fmt.Errorf("error while exec register query: %w", err)
	}

	return domain.User{UUID: userUUID}, nil
}
