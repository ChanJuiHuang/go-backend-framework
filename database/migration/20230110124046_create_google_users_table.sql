-- +goose Up
-- +goose StatementBegin
CREATE TABLE google_users (
  google_id VARCHAR(255) NOT NULL PRIMARY KEY,
  user_id BIGINT UNSIGNED UNIQUE NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS google_users;
-- +goose StatementEnd
