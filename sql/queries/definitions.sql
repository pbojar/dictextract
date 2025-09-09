-- name: CreateDefinition :one
INSERT INTO definitions (id, word_id, pos_id, "definition")
VALUES (
    DEFAULT,
    $1,
    $2,
    $3
)
RETURNING *;
