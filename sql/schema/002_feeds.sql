-- +goose Up
CREATE TABLE feeds (
  id UUID PRIMARY KEY,
  url text UNIQUE NOT NULL,
  user_id UUID NOT NULL,
  name text NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  FOREIGN KEY(user_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
