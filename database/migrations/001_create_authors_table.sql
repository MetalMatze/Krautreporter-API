-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE authors (
    id integer PRIMARY KEY NOT NULL,
    ordering integer,
    name text,
    title text,
    url text,
    biography text,
    social_media text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE authors;
