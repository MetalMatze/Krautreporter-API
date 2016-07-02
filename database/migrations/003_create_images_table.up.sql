CREATE TABLE images (
    id SERIAL,
    width integer,
    src text,
    imageable_id integer,
    imageable_type text
);
