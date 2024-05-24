-- +goose Up
-- +goose StatementBegin
CREATE TABLE http_apis (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  method VARCHAR(255) NOT NULL,
  path VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  UNIQUE(method, path) ON CONFLICT REPLACE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS http_apis;
-- +goose StatementEnd
