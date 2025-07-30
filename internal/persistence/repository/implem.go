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

func NewVersionRepository(queries *db.Queries) VersionRepository {
	return &versionRepository{queries: queries}
}

func NewDataModelRepository(queries *db.Queries) DatamodelRepository {
	return &dataModelRepository{queries: queries}
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
