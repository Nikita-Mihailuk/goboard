package postgres

import (
	"context"
	"errors"
	"github.com/Nikita-Mihailuk/goboard/backend/user_service/internal/domain/dto"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) SaveUser(ctx context.Context, input dto.CreateUserInput) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	queryInsertUser := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)`
	_, err = tx.Exec(ctx, queryInsertUser, input.Name, input.Email, input.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return ErrUserExists
		}
		return err
	}

	return tx.Commit(ctx)
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (dto.LoginUserOutput, error) {
	var outputUser dto.LoginUserOutput
	queryFindUserByEmail := `SELECT id, role, password_hash FROM users WHERE email = $1`
	err := s.db.QueryRow(ctx, queryFindUserByEmail, email).Scan(&outputUser.ID, &outputUser.Role, &outputUser.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.LoginUserOutput{}, ErrUserNotFound
		}
		return dto.LoginUserOutput{}, err
	}
	return outputUser, nil
}

func (s *Storage) GetUserByID(ctx context.Context, id int64) (dto.GetUserByIDOutput, error) {
	var outputUser dto.GetUserByIDOutput
	queryFindUserByID := `SELECT email, photo_url, name FROM users WHERE id = $1`
	err := s.db.QueryRow(ctx, queryFindUserByID, id).Scan(&outputUser.Email, &outputUser.PhotoURL, &outputUser.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.GetUserByIDOutput{}, ErrUserNotFound
		}
		return dto.GetUserByIDOutput{}, err
	}
	return outputUser, nil
}

func (s *Storage) GetUserUpdateByID(ctx context.Context, id int64) (dto.UserForUpdate, error) {
	var outputUser dto.UserForUpdate
	queryFindUserByID := `SELECT  password_hash ,photo_url, name FROM users WHERE id = $1`
	err := s.db.QueryRow(ctx, queryFindUserByID, id).Scan(&outputUser.PasswordHash, &outputUser.PhotoUrl, &outputUser.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.UserForUpdate{}, ErrUserNotFound
		}
		return dto.UserForUpdate{}, err
	}
	outputUser.ID = id
	return outputUser, nil
}

func (s *Storage) RefreshUser(ctx context.Context, input dto.UserForUpdate) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	queryUpdateUser := `UPDATE users SET password_hash = $1, name = $2, photo_url = $3 WHERE id = $4`
	_, err = tx.Exec(ctx, queryUpdateUser, input.PasswordHash, input.Name, input.PhotoUrl.String, input.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
