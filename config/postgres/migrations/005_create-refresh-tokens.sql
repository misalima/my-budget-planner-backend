-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS refresh_tokens
(
    "ID" uuid NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    token character varying(512) NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    PRIMARY KEY ("ID"),
    FOREIGN KEY (user_id) REFERENCES users ("ID")
);
---- create above / drop below ----

DROP TABLE IF EXISTS refresh_tokens;

