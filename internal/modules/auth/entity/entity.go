package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID           string       `db:"id"`
	Email        string       `db:"email"`
	PasswordHash string       `db:"password_hash"`
	Status       string       `db:"status"`
	CreatedAt    time.Time    `db:"created_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

func (u *User) IsActive() bool {
	return u.Status == "active"
}

type Role struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type RefreshToken struct {
	ID        string       `db:"id"`
	UserID    string       `db:"user_id"`
	TokenHash string       `db:"token_hash"`
	ExpiresAt time.Time    `db:"expires_at"`
	RevokedAt sql.NullTime `db:"revoked_at"`
	CreatedAt time.Time    `db:"created_at"`
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt.Valid
}
