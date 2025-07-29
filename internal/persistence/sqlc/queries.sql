-- name: CreateDatamodel :one
INSERT INTO datamodel (name)
VALUES ($1)
    RETURNING id, name;

-- name: GetDatamodel :one
SELECT json
FROM version
WHERE object_id = 123 AND object_type = 'datamodel'
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

