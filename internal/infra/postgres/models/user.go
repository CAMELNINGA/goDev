package models

import "Yaratam/internal/domain"

type User struct {
	ID       int    `db:"id"`
	UserName string `db:"username"`
	ChatID   int    `db:"chat_id"`
	PathID   int    `db:"path_id"`
}

func (u *User) ToDomain() *domain.User {
	return &domain.User{
		UserName: u.UserName,
		ChatID:   u.ChatID,
		PathID:   u.PathID,
	}
}
