-- name: CreateDatamodel :one
INSERT INTO datamodel (name)
VALUES ($1)
    RETURNING id, name;

-- name: GetDatamodel :one
SELECT json
FROM version
WHERE object_id = $1 AND object_type = 'datamodel'
ORDER BY version DESC
LIMIT 1;

-- name: ListDatamodels :many
SELECT id, name
FROM datamodel
ORDER BY id;

-- name: UpdateDatamodel :one
UPDATE datamodel
SET name = $2
WHERE id = $1
    RETURNING id, name;

-- name: DeleteDatamodel :exec
DELETE FROM datamodel
WHERE id = $1;

-- name: CreateVersion :one
INSERT INTO version (object_type, object_id, json, version, action, actor, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, object_type, object_id, json, version, action, actor, created_at;

-- name: GetVersionByID :one
SELECT id, object_type, object_id, json, version, action, actor, created_at
FROM version
WHERE id = $1;

-- name: ListVersionsByObject :many
SELECT id, object_type, object_id, json, version, action, actor, created_at
FROM version
WHERE object_type = $1 AND object_id = $2
ORDER BY version DESC;

