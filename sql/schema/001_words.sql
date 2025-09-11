-- +goose Up
CREATE TABLE words(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    word TEXT UNIQUE NOT NULL,
    CONSTRAINT no_whitespace CHECK (word !~ '\s'),
    CONSTRAINT word_lowercase_ck CHECK (word = LOWER(word))
);

-- +goose Down
DROP TABLE words;
