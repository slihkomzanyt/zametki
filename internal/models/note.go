package models

import "time"

type Note struct {
	ID        int64     ` json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// DTO для запроса создания заметки
type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// DTO для запроса обновления заметки
type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
