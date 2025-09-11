-- name: CreateWord :one
INSERT INTO words (id, word)
VALUES (
    DEFAULT, 
    $1
) 
RETURNING *;

-- name: GetIDByWord :one
SELECT id FROM words WHERE word=$1;
