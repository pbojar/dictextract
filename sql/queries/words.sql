-- name: CreateWord :one
INSERT INTO words (id, word)
VALUES (
    DEFAULT, 
    $1
) 
RETURNING *;
