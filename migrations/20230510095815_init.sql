-- noinspection SqlNoDataSourceInspectionForFile

-- +goose Up
-- +goose StatementBegin
CREATE TABLE url_shorten(
  id BIGSERIAL PRIMARY KEY,
  token varchar(30) NOT NULL,
  url_shorten varchar(1000) NOT NULL,
  url varchar(4000) NOT NULL,
  created_at timestamp NOT NULL DEFAULT NOW(),
  expired_at timestamp NOT NULL
);

CREATE TABLE users_authentication(
    id BIGSERIAL PRIMARY KEY,
    user_name varchar(30),
    password varchar(100),
    email varchar(100)
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE url_shorten;
DROP TABLE users_authentication;
-- +goose StatementEnd
