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
