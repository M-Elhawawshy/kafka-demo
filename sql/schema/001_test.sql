-- +goose Up
CREATE TABLE test(
    id UUID PRIMARY KEY,
    content TEXT
);

-- +goose Down
DROP TABLE test;