-- +goose Up
CREATE TABLE parts_of_speech(
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pos TEXT UNIQUE NOT NULL,
    CONSTRAINT no_whitespace CHECK (pos !~ '\s'),
    CONSTRAINT pos_lowercase_ck CHECK (pos = LOWER(pos))
);

-- +goose Down
DROP TABLE parts_of_speech;
