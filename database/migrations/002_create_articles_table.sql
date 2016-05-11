-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE articles (
    id integer PRIMARY KEY NOT NULL,
    ordering integer,
    title text,
    headline text,
    date text,
    preview boolean,
    url text,
    excerpt text,
    content text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    author_id integer
);

--CREATE INDEX author_id_idx ON authors (author_id);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE articles;
