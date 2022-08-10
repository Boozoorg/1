package client

import (
	"context"
	"database/sql"
	"errors"
	"log"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var ErrNotFound = errors.New("item not fount")

var ErrInternal = errors.New("internal error")

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

type Account struct {
	Name string `json:"name"`
}

func (s *Service) Registration(ctx context.Context, name string) (*Account, error) {
	item := &Account{}

	err := s.db.QueryRowContext(ctx, `
		SELECT * FROM client
	`, name).Scan(&item.Name)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}
