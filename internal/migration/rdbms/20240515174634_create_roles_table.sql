-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL,
  is_public BOOLEAN NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
