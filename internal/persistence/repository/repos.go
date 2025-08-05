package repository

import (
	"cognyx/psychic-robot/persistence/db"
	"context"
)

// Interface pour User
type UserRepository interface {
	Create(ctx context.Context, user db.User) (db.User, error)
	GetByID(ctx context.Context, id string) (db.User, error)
	GetByEmail(ctx context.Context, email string) (db.User, error)
	List(ctx context.Context, limit, offset int32) ([]db.User, error)
	Update(ctx context.Context, user db.User) (db.User, error)
	Delete(ctx context.Context, id string) error
}
