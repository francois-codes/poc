package repository

import (
	"cognyx/psychic-robot/persistence/db"
	"context"
)

type PostgresUserRepository struct {
	q *db.Queries
}

func NewUserRepository(q *db.Queries) *PostgresUserRepository {
	return &PostgresUserRepository{q: q}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user db.User) (db.User, error) {
	u, err := r.q.CreateUser(ctx, db.CreateUserParams{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Roles: user.Roles,
	})
	if err != nil {
		return db.User{}, err
	}
	return u, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (db.User, error) {
	u, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return db.User{}, err
	}
	return u, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (db.User, error) {
	u, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, err
	}
	return u, nil
}

func (r *PostgresUserRepository) List(ctx context.Context, limit, offset int32) ([]db.User, error) {
	users, err := r.q.ListUsers(ctx, db.ListUsersParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user db.User) (db.User, error) {
	u, err := r.q.UpdateUser(ctx, db.UpdateUserParams{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Roles: user.Roles,
	})
	if err != nil {
		return db.User{}, err
	}
	return u, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	return r.q.DeleteUser(ctx, id)
}
