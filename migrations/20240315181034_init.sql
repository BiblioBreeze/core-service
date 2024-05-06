-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id         bigserial PRIMARY KEY,
    email      varchar(255),
    first_name varchar(255),
    last_name  varchar(255),
    password   varchar(255)
);

CREATE TABLE books
(
    id                 bigserial PRIMARY KEY,
    belongs_to_user_id bigint REFERENCES users (id),
    name               varchar(255),
    author             varchar(255),
    genre              varchar(255),
    description        varchar(1000),
    latitude           decimal(8, 6),
    longitude          decimal(9, 6)
);

CREATE TABLE exchange_requests
(
    id           bigserial PRIMARY KEY,
    from_user_id bigint REFERENCES users (id),
    title        varchar(255),
    condition    varchar(1000)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE exchange_requests;
DROP TABLE books;
DROP TABLE users;
-- +goose StatementEnd
