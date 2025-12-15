-- public.urls definition

-- Drop table

-- DROP TABLE public.urls;

CREATE TABLE public.urls
(
    id         bigserial NOT NULL,
    urls             text      NOT NULL,
    short_urls       text      NOT NULL,
    click_count      bigint    DEFAULT 0,
    last_accessed_at timestamptz    NULL,
    created_at       timestamptz    NOT NULL,
    updated_at timestamptz    NOT NULL,
    deleted_at timestamptz    NULL
)
    PARTITION BY LIST (deleted_at);

CREATE TABLE public.url__active PARTITION OF public.urls FOR VALUES IN (NULL);
CREATE TABLE public.url_inactive PARTITION OF public.urls DEFAULT;

CREATE INDEX urls_id_idx ON ONLY public.urls USING btree (id, deleted_at);
CREATE INDEX urls_short_urls_idx ON ONLY public.urls USING btree (short_urls, deleted_at);
CREATE INDEX urls_urls_idx ON ONLY public.urls USING btree (urls, deleted_at);

