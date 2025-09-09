-- +goose Up
CREATE TABLE definitions(
    id SERIAL PRIMARY KEY,
    word_id SERIAL NOT NULL,
    CONSTRAINT fk_word_id
    FOREIGN KEY (word_id)
    REFERENCES words(id)
    ON DELETE CASCADE,
    pos_id SERIAL NOT NULL,
    CONSTRAINT fk_pos_id
    FOREIGN KEY (pos_id)
    REFERENCES parts_of_speech(id)
    ON DELETE CASCADE,
    "definition" TEXT NOT NULL,
    CONSTRAINT uc_definition
    UNIQUE (word_id, pos_id)
);

-- +goose Down
DROP TABLE definitions;
