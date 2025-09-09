-- name: CreatePos :one
INSERT INTO parts_of_speech (id, pos)
VALUES (
    DEFAULT, 
    $1
) 
RETURNING *;
