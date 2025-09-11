-- name: CreatePos :one
INSERT INTO parts_of_speech (id, pos)
VALUES (
    DEFAULT, 
    $1
) 
RETURNING *;

-- name: GetIDByPos :one
SELECT id FROM parts_of_speech WHERE pos=$1;
