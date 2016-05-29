-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE crawls (
    id SERIAL,
    next timestamp with time zone,
    crawlable_id integer,
    crawlable_type text
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE crawls;
