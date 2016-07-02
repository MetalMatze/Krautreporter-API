CREATE TABLE crawls (
    id SERIAL,
    next timestamp with time zone,
    crawlable_id integer,
    crawlable_type text
);
