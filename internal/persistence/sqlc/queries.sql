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
VALUES ($1, $2, $3, $4, $5, $6, NOW())
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


-- name: SearchUsers :many
SELECT * FROM users
WHERE (@email::text IS NULL OR email = @email)
  AND (@status::text IS NULL OR status = @status)
  AND (@role::text IS NULL OR role = @role)
ORDER BY
    CASE
        WHEN @orderby::text = 'email' THEN email
        WHEN @orderby::text = 'created_at' THEN created_at
        WHEN @orderby::text = 'role' THEN role
        ELSE id
        END
LIMIT $1 OFFSET $2;

-- name: FilterUsers :many
SELECT * FROM users
WHERE (CASE WHEN @is_email::bool THEN email = @email ELSE TRUE END)
  AND (CASE WHEN @like_email::bool THEN email LIKE @email ELSE TRUE END)
  AND (CASE WHEN @status::text IS NOT NULL THEN status = @status ELSE TRUE END)
  AND (CASE WHEN @role::text IS NOT NULL THEN role = @role ELSE TRUE END)
ORDER BY
    CASE WHEN @email_asc::bool THEN email END asc,
    CASE WHEN @email_desc::bool THEN email END desc,
    CASE WHEN @status_asc::bool THEN status END asc,
    CASE WHEN @status_desc::bool THEN status END desc
LIMIT $1 OFFSET $2;

-- CRUD Operations for Users

-- name: CreateUser :one
INSERT INTO users (email, status, role)
VALUES ($1, $2, $3)
RETURNING id, email, status, role, created_at;

-- name: GetUser :one
SELECT id, email, status, role, created_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, status, role, created_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, email, status, role, created_at
FROM users
ORDER BY created_at DESC;

-- name: UpdateUser :one
UPDATE users
SET email = $2, status = $3, role = $4
WHERE id = $1
RETURNING id, email, status, role, created_at;

-- name: UpdateUserStatus :one
UPDATE users
SET status = $2
WHERE id = $1
RETURNING id, email, status, role, created_at;

-- name: UpdateUserRole :one
UPDATE users
SET role = $2
WHERE id = $1
RETURNING id, email, status, role, created_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
