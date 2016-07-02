CREATE TABLE articles (
    id integer PRIMARY KEY NOT NULL,
    ordering integer,
    title text,
    headline text,
    date timestamp with time zone,
    preview boolean,
    url text,
    excerpt text,
    content text,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    author_id integer
);
