package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type (
	// users table
	User struct {
		ID        uuid.UUID `db:"id"`
		Name      string    `db:"name"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
	}
)

func randomDefaultName() string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return "Player-" + string(b)
}

func (r *Repository) CreateUser(ctx context.Context) (uuid.UUID, error) {
	userID := uuid.New()
	name := randomDefaultName()
	if _, err := r.db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (?, ?)", userID, name); err != nil {
		return uuid.Nil, fmt.Errorf("insert user: %w", err)
	}

	return userID, nil
}

func (r *Repository) UpdateUserName(ctx context.Context, userID uuid.UUID, name string) error {
	res, err := r.db.ExecContext(ctx, "UPDATE users SET name = ? WHERE id = ?", name, userID)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *Repository) GetUser(ctx context.Context, userID uuid.UUID) (*User, error) {
	user := &User{}
	if err := r.db.GetContext(ctx, user, "SELECT * FROM users WHERE id = ?", userID); err != nil {
		return nil, fmt.Errorf("select user: %w", err)
	}
	return user, nil
}
