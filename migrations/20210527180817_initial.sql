-- Add migration script here
-- Your SQL goes here
CREATE TABLE users (
  id BIGINT UNIQUE PRIMARY KEY,
  name text,
  automatic_sending bool
--  created_at timestamp,
--  updated_at timestamp,
--  deleted_at timestamp
);

-- Your SQL goes here
CREATE TABLE keywords (
    id SERIAL UNIQUE PRIMARY KEY,
    searchterm text,
    user_id BIGINT NOT NULL REFERENCES users (id)
);