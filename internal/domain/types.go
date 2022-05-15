package domain

type ContextKey string

const ContextUserID ContextKey = "ctx_user_id"

type User struct {
	UserName string
	ChatID   int
}
