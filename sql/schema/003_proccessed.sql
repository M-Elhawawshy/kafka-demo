-- +goose Up
CREATE TABLE processed (
    id UUID PRIMARY KEY,
    content TEXT
);

-- +goose Down
DROP TABLE processed;