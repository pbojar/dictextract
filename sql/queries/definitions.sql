-- name: CreateDefinition :one
INSERT INTO definitions (id, word_id, pos_id, "definition")
VALUES (
    DEFAULT,
    $1,
    $2,
    $3
)
RETURNING *;

-- name: DefinitionExists :one
SELECT EXISTS(SELECT 1 FROM definitions WHERE word_id=$1 AND pos_id=$2);
