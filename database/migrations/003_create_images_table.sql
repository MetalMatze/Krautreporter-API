-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE images (
    id SERIAL,
    width integer,
    src text,
    imageable_id integer,
    imageable_type text
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE images;
