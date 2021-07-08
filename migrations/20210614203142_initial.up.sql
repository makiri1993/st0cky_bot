-- Add up migration script here

CREATE TABLE users (
  id BIGINT UNIQUE PRIMARY KEY,
  name text,
  automatic_sending bool
--  created_at timestamp,
--  updated_at timestamp,
--  deleted_at timestamp
);

CREATE TABLE keywords (
    id SERIAL UNIQUE PRIMARY KEY,
    searchterm text,
    user_id BIGINT NOT NULL REFERENCES users (id)
);

CREATE TABLE news (
    id numeric UNIQUE PRIMARY KEY,
    title text,
    url text,
    description text,
    date_published TIMESTAMP,
    sent bool,
    user_id BIGINT NOT NULL REFERENCES users (id)
);