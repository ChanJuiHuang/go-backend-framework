-- +goose Up
-- +goose StatementBegin
CREATE TABLE events(
    name String,
    created_at DateTime
)
ENGINE = MergeTree
PRIMARY KEY (created_at, name)
PARTITION BY toYYYYMM(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
