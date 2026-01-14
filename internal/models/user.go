package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// DTO для запроса создания пользователя
type CreateUserRequest struct {
	Username string `json:"username"`
}
