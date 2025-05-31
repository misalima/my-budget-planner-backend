CREATE TABLE IF NOT EXISTS budgets
(
    "ID" serial PRIMARY KEY NOT NULL,
    user_id uuid NOT NULL,
    budget_name character varying(255) NOT NULL,
    description character varying(255),
    amount double precision NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users ("ID") ON DELETE CASCADE
);
---- create above / drop below ----

DROP TABLE IF EXISTS budgets;