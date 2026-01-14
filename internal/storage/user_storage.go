package storage

import (
	"context"
	"database/sql"
	"errors"

	"zametki/internal/models"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

type UserStorage interface {
	CreateUser(ctx context.Context, username string) (models.User, error)
	GetUserByID(ctx context.Context, id int64) (models.User, error)
}

func (p *Postgres) CreateUser(ctx context.Context, username string) (models.User, error) {
	const q = `
  INSERT INTO users (username)
  VALUES ($1)
  RETURNING id, username, created_at
  `

	var u models.User
	err := p.DB.QueryRowContext(ctx, q, username).Scan(&u.ID, &u.Username, &u.CreatedAt)
	if err != nil {
		// UNIQUE violation удобнее обрабатывать позже в handler,
		// но можно и здесь. Пока оставим как есть.
		return models.User{}, err
	}
	return u, nil
}

func (p *Postgres) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	const q = `
  SELECT id, username, created_at
  FROM users
  WHERE id = $1
  `

	var u models.User
	err := p.DB.QueryRowContext(ctx, q, id).Scan(&u.ID, &u.Username, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, ErrNotFound
		}
		return models.User{}, err
	}
	return u, nil
}
