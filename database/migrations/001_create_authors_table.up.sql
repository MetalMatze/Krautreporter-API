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
