package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	trmsqlx "github.com/avito-tech/go-transaction-manager/drivers/sqlx/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/inna-maikut/avito-pvz/internal/model"
)

type UserRepository struct {
	db     *sqlx.DB
	getter *trmsqlx.CtxGetter
}

func NewUserRepository(db *sqlx.DB, getter *trmsqlx.CtxGetter) (*UserRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	if getter == nil {
		return nil, errors.New("getter is nil")
	}

	return &UserRepository{
		db:     db,
		getter: getter,
	}, nil
}

func (r *UserRepository) trOrDB(ctx context.Context) trmsqlx.Tr {
	return r.getter.DefaultTrOrDB(ctx, r.db)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user User

	q := "SELECT id, email, password, user_role FROM users WHERE email = $1"

	err := r.trOrDB(ctx).GetContext(ctx, &user, q, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("db.GetContext: %w", err)
	}

	return &model.User{
		UserID:   model.UserID(user.ID),
		Email:    user.Email,
		Password: user.Password,
		UserRole: model.UserRole(user.Role),
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash string, role model.UserRole) (*model.User, error) {
	q := "INSERT INTO users (email, password, user_role) values " +
		"($1, $2, $3) " + // use binding to avoid SQL injection
		"ON CONFLICT DO NOTHING " +
		"RETURNING ID"

	var ID uuid.UUID
	err := r.trOrDB(ctx).GetContext(ctx, &ID, q, email, passwordHash, role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("db.GetContext: %w", err)
	}

	return &model.User{
		UserID:   model.UserID(ID),
		Email:    email,
		Password: passwordHash,
		UserRole: role,
	}, nil
}
