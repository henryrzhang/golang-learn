-- name: CreateDrama :one
INSERT INTO drama_info (
    drama_no, title, outline, cover_image, characters,
    character_relation_desc, status, task_no, create_by, update_by,
    create_at, update_at, deleted
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: GetDramaByID :one
SELECT * FROM drama_info WHERE id = $1 AND deleted = FALSE;

-- name: GetDramaByNo :one
SELECT * FROM drama_info WHERE drama_no = $1 AND deleted = FALSE;

-- name: ListDramas :many
SELECT * FROM drama_info
WHERE deleted = FALSE
ORDER BY id DESC
LIMIT $1 OFFSET $2;

-- name: UpdateDrama :one
UPDATE drama_info
SET
    title = COALESCE(sqlc.narg('title'), title),
    outline = COALESCE(sqlc.narg('outline'), outline),
    cover_image = COALESCE(sqlc.narg('cover_image'), cover_image),
    characters = COALESCE(sqlc.narg('characters'), characters),
    character_relation_desc = COALESCE(sqlc.narg('character_relation_desc'), character_relation_desc),
    status = COALESCE(sqlc.narg('status'), status),
    task_no = COALESCE(sqlc.narg('task_no'), task_no),
    update_by = COALESCE(sqlc.narg('update_by'), update_by),
    update_at = sqlc.arg('update_at')
WHERE id = sqlc.arg('id') AND deleted = FALSE
RETURNING *;

-- name: SoftDeleteDrama :one
UPDATE drama_info
SET deleted = TRUE, update_at = $2
WHERE id = $1 AND deleted = FALSE
RETURNING *;
