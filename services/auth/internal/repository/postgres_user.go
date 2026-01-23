package repository

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	// #region agent log
	if f, err := os.OpenFile("/home/smart/Belajar/project/event-driven/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
		json.NewEncoder(f).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "B", "location": "postgres_user.go:17", "message": "FindByUsername called", "data": map[string]interface{}{"username": username}, "timestamp": os.Getpid()})
		f.Close()
	}
	// #endregion

	var user User

	err := r.db.QueryRow(ctx,
		`SELECT id, username, password FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password)

	// #region agent log
	if f, err2 := os.OpenFile("/home/smart/Belajar/project/event-driven/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		isNoRows := errors.Is(err, pgx.ErrNoRows)
		json.NewEncoder(f).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "B", "location": "postgres_user.go:25", "message": "QueryRow result", "data": map[string]interface{}{"error": err != nil, "isNoRows": isNoRows, "errorType": func() string {
			if err != nil {
				return err.Error()
			} else {
				return "nil"
			}
		}()}, "timestamp": os.Getpid()})
		f.Close()
	}
	// #endregion

	if err != nil {
		// #region agent log
		if f, err2 := os.OpenFile("/home/smart/Belajar/project/event-driven/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
			json.NewEncoder(f).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "C", "location": "postgres_user.go:30", "message": "Returning error", "data": map[string]interface{}{"error": err.Error(), "isNoRows": errors.Is(err, pgx.ErrNoRows)}, "timestamp": os.Getpid()})
			f.Close()
		}
		// #endregion
		return nil, err
	}

	// #region agent log
	if f, err2 := os.OpenFile("/home/smart/Belajar/project/event-driven/.cursor/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err2 == nil {
		json.NewEncoder(f).Encode(map[string]interface{}{"sessionId": "debug-session", "runId": "run1", "hypothesisId": "B", "location": "postgres_user.go:36", "message": "User found successfully", "data": map[string]interface{}{"userId": user.ID, "username": user.Username}, "timestamp": os.Getpid()})
		f.Close()
	}
	// #endregion

	return &user, nil
}
