package repository

import (
	"cognyx/psychic-robot/persistence/db"
	"context"
	"encoding/json"
)

// Interface pour Datamodel
type DatamodelRepository interface {
	Create(ctx context.Context, name string) (db.Datamodel, error)
	Update(ctx context.Context, id int64, name string) (db.Datamodel, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]db.Datamodel, error)
}

// Interface pour Version
type VersionRepository interface {
	Create(ctx context.Context, v db.CreateVersionParams) (db.Version, error)
	GetLatestByDatamodelID(ctx context.Context, id int64) (json.RawMessage, error)
	ListByObject(ctx context.Context, objectType string, objectID int64) ([]db.Version, error)
	GetByID(ctx context.Context, id int64) (db.Version, error)
}

// Interface pour User
type UserRepository interface {
	Create(ctx context.Context, params db.CreateUserParams) (db.User, error)
	GetByID(ctx context.Context, id int64) (db.User, error)
	GetByEmail(ctx context.Context, email string) (db.User, error)
	List(ctx context.Context) ([]db.User, error)
	Update(ctx context.Context, params db.UpdateUserParams) (db.User, error)
	UpdateStatus(ctx context.Context, params db.UpdateUserStatusParams) (db.User, error)
	UpdateRole(ctx context.Context, params db.UpdateUserRoleParams) (db.User, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int64, error)
	Search(ctx context.Context, params db.SearchUsersParams) ([]db.User, error)
	Filter(ctx context.Context, params db.FilterUsersParams) ([]db.User, error)
}
