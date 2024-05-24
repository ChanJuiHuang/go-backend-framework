-- +goose Up
-- +goose StatementBegin
CREATE TABLE http_apis (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  method VARCHAR(255) NOT NULL,
  path VARCHAR(255) NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  UNIQUE method_path_uniq(method, path),
  INDEX path_idx(path)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS http_apis;
-- +goose StatementEnd
