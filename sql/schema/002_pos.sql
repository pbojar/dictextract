-- +goose Up
CREATE TABLE parts_of_speech(
    id SERIAL PRIMARY KEY,
    pos TEXT UNIQUE NOT NULL,
    CONSTRAINT no_whitespace CHECK (pos !~ '\s')
);

-- +goose Down
DROP TABLE parts_of_speech;
