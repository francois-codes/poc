package repository

import (
	"cognyx/psychic-robot/persistence/db"
	"context"
	"encoding/json"
)

type versionRepository struct {
	queries *db.Queries
}
type dataModelRepository struct {
	queries *db.Queries
}

type userRepository struct {
	queries *db.Queries
}

func NewVersionRepository(queries *db.Queries) VersionRepository {
	return &versionRepository{queries: queries}
	
}

func NewDataModelRepository(queries *db.Queries) DatamodelRepository {
	return &dataModelRepository{queries: queries}
}

func NewUserRepository(queries *db.Queries) UserRepository {
	return &userRepository{queries: queries}
}

// --- Datamodel ---

func (r *dataModelRepository) Create(ctx context.Context, name string) (db.Datamodel, error) {
	return r.queries.CreateDatamodel(ctx, name)
}

func (r *dataModelRepository) Update(ctx context.Context, id int64, name string) (db.Datamodel, error) {
	return r.queries.UpdateDatamodel(ctx, db.UpdateDatamodelParams{ID: id, Name: name})
}

func (r *dataModelRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteDatamodel(ctx, id)
}

func (r *dataModelRepository) List(ctx context.Context) ([]db.Datamodel, error) {
	return r.queries.ListDatamodels(ctx)
}

// --- Version ---

func (r *versionRepository) Create(ctx context.Context, v db.CreateVersionParams) (db.Version, error) {
	return r.queries.CreateVersion(ctx, v)
}

func (r *versionRepository) GetLatestByDatamodelID(ctx context.Context, id int64) (json.RawMessage, error) {
	return r.queries.GetDatamodel(ctx, id)
}

func (r *versionRepository) ListByObject(ctx context.Context, objectType string, objectID int64) ([]db.Version, error) {
	return r.queries.ListVersionsByObject(ctx, db.ListVersionsByObjectParams{
		ObjectType: objectType,
		ObjectID:   objectID,
	})
}

func (r *versionRepository) GetByID(ctx context.Context, id int64) (db.Version, error) {
	return r.queries.GetVersionByID(ctx, id)
}

// --- User ---

func (r *userRepository) Create(ctx context.Context, params db.CreateUserParams) (db.User, error) {
	return r.queries.CreateUser(ctx, params)
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (db.User, error) {
	return r.queries.GetUser(ctx, id)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (db.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}

func (r *userRepository) List(ctx context.Context) ([]db.User, error) {
	return r.queries.ListUsers(ctx)
}

func (r *userRepository) Update(ctx context.Context, params db.UpdateUserParams) (db.User, error) {
	return r.queries.UpdateUser(ctx, params)
}

func (r *userRepository) UpdateStatus(ctx context.Context, params db.UpdateUserStatusParams) (db.User, error) {
	return r.queries.UpdateUserStatus(ctx, params)
}

func (r *userRepository) UpdateRole(ctx context.Context, params db.UpdateUserRoleParams) (db.User, error) {
	return r.queries.UpdateUserRole(ctx, params)
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteUser(ctx, id)
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}

func (r *userRepository) Search(ctx context.Context, params db.SearchUsersParams) ([]db.User, error) {
	return r.queries.SearchUsers(ctx, params)
}

func (r *userRepository) Filter(ctx context.Context, params db.FilterUsersParams) ([]db.User, error) {
	return r.queries.FilterUsers(ctx, params)
}
