package storage

import (
	"context"
	"database/sql"
	"errors"

	"zametki/internal/models"
)

type NoteStorage interface {
	CreateNote(ctx context.Context, userID int64, title, content string) (models.Note, error)
	GetNotesByUser(
		ctx context.Context,
		userID int64,
		limit, offset int,
		sortAsc bool,
	) ([]models.Note, error)
	GetNoteByID(ctx context.Context, userID, noteID int64) (models.Note, error)
	UpdateNote(ctx context.Context, userID, noteID int64, title, content string) (models.Note, error)
	DeleteNote(ctx context.Context, userID, noteID int64) error
}

func (p *Postgres) CreateNote(
	ctx context.Context,
	userID int64,
	title, content string,
) (models.Note, error) {
	const q = `
  INSERT INTO notes (user_id, title, content)
  VALUES ($1, $2, $3)
  RETURNING id, user_id, title, content, created_at, updated_at
 `

	var n models.Note
	err := p.DB.QueryRowContext(ctx, q, userID, title, content).
		Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		return models.Note{}, err
	}
	return n, nil
}

func (p *Postgres) GetNotesByUser(
	ctx context.Context,
	userID int64,
	limit, offset int,
	sortAsc bool,
) ([]models.Note, error) {
	order := "DESC"
	if sortAsc {
		order = "ASC"
	}

	q := `
  SELECT id, user_id, title, content, created_at, updated_at
  FROM notes
  WHERE user_id = $1
  ORDER BY created_at ` + order + `
  LIMIT $2 OFFSET $3
 `

	rows, err := p.DB.QueryContext(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Title,
			&n.Content,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, nil
}

func (p *Postgres) GetNoteByID(
	ctx context.Context,
	userID, noteID int64,
) (models.Note, error) {
	const q = `
  SELECT id, user_id, title, content, created_at, updated_at
  FROM notes
  WHERE id = $1 AND user_id = $2
 `

	var n models.Note
	err := p.DB.QueryRowContext(ctx, q, noteID, userID).
		Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Note{}, ErrNotFound
		}
		return models.Note{}, err
	}
	return n, nil
}

func (p *Postgres) UpdateNote(
	ctx context.Context,
	userID, noteID int64,
	title, content string,
) (models.Note, error) {
	const q = `
  UPDATE notes
  SET title = $1,
      content = $2,
      updated_at = now()
  WHERE id = $3 AND user_id = $4
  RETURNING id, user_id, title, content, created_at, updated_at
 `

	var n models.Note
	err := p.DB.QueryRowContext(ctx, q, title, content, noteID, userID).
		Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Note{}, ErrNotFound
		}
		return models.Note{}, err
	}
	return n, nil
}

func (p *Postgres) DeleteNote(
	ctx context.Context,
	userID, noteID int64,
) error {
	const q = `
  DELETE FROM notes
  WHERE id = $1 AND user_id = $2
 `

	res, err := p.DB.ExecContext(ctx, q, noteID, userID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
