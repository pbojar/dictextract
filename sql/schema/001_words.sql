-- +goose Up
CREATE TABLE words(
    id SERIAL PRIMARY KEY,
    word TEXT UNIQUE NOT NULL,
    CONSTRAINT no_whitespace CHECK (word !~ '\s')
);

-- +goose Down
DROP TABLE words;
