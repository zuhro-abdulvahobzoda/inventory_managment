package user

import (
	"context"
	"errors"
	"fmt"
	"inventory-api/internal/models"
	"inventory-api/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueViolationCode = "23505"
	// singleAdminIndex is the unique index (in migrations) that allows only one admin.
	singleAdminIndex = "users_single_admin"
	userColumns      = "id, username, password_hash, role"
)

var _ repository.UserRepository = (*userRepo)(nil)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *userRepo {
	return &userRepo{pool: pool}
}

func (r *userRepo) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (username, password_hash, role)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		user.Username, user.PasswordHash, user.Role,
	).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			if pgErr.ConstraintName == singleAdminIndex {
				return 0, models.ErrAdminAlreadyExists
			}
			return 0, models.ErrUsernameAlreadyExists
		}
		return 0, fmt.Errorf("create user: %w", err)
	}

	return user.ID, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user, err := r.getOne(ctx, `SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := r.getOne(ctx, `SELECT `+userColumns+` FROM users WHERE username = $1`, username)
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	return user, nil
}

func (r *userRepo) getOne(ctx context.Context, query string, args ...any) (*models.User, error) {
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[dbModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	user := m.toModel()
	return &user, nil
}
