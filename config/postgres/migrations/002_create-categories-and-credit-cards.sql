CREATE TABLE IF NOT EXISTS credit_cards
(
    "ID" uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    card_name character varying NOT NULL,
    total_limit double precision NOT NULL,
    current_limit double precision,
    due_date integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users ("ID") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories
(
    "ID" serial PRIMARY KEY NOT NULL,
    category_name character varying NOT NULL,
    user_id uuid DEFAULT NULL,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users ("ID") ON DELETE CASCADE
);

INSERT INTO categories (category_name)
VALUES ('Groceries'),
       ('Food'),
       ('Transportation'),
       ('Health'),
       ('Education'),
       ('Entertainment'),
       ('Others');

---- create above / drop below ----

DROP TABLE IF EXISTS credit_cards;
DROP TABLE IF EXISTS categories;
