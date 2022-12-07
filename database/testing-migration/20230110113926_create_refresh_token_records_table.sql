-- +goose Up
-- +goose StatementBegin
CREATE TABLE refresh_token_records (
  refresh_token VARCHAR(510) PRIMARY KEY,
  user_id BIGINT UNSIGNED NOT NULL NOT NULL,
  device VARCHAR(255) CHECK( device IN ('web', 'ios', 'android')) NOT NULL,
  expire_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS refresh_token_records;
-- +goose StatementEnd
