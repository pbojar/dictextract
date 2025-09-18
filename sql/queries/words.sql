-- name: CreateWord :one
INSERT INTO words (id, word)
VALUES (
    DEFAULT, 
    $1
) 
RETURNING *;

-- name: GetIDByWord :one
SELECT id FROM words WHERE word=$1;

-- name: GetWordsWithLenInRangeSorted :many
SELECT word FROM words WHERE CHAR_LENGTH(word) BETWEEN @minLen AND @maxLen ORDER BY word ASC;
