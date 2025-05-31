CREATE TABLE IF NOT EXISTS users
(
    "ID" uuid NOT NULL DEFAULT gen_random_uuid(),
    email character varying(255) NOT NULL,
    username character varying(255) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    password_hash character(255) NOT NULL,
    profile_picture character varying(255) NOT NULL DEFAULT '',
    income double precision NOT NULL DEFAULT 0,
    expenditure_limit double precision NOT NULL DEFAULT 0,
    created_at timestamp with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp with time zone NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("ID")
);
---- create above / drop below ----

DROP TABLE IF EXISTS users;