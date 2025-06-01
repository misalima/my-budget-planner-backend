CREATE TABLE IF NOT EXISTS simple_expense (
    ID uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    category_id int NOT NULL,
    amount double precision NOT NULL,
    description character varying(255),
    date date NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (ID) ON DELETE CASCADE,
    CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES categories (ID) ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS recurring_expense (
    ID uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    category_id int NOT NULL,
    amount double precision NOT NULL,
    description character varying(255),
    date date NOT NULL,
    card_id uuid,
    start_date date NOT NULL,
    end_date date,
    frequency character varying NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (ID) ON DELETE CASCADE,
    CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES categories (ID) ON DELETE RESTRICT,
    CONSTRAINT fk_card_id FOREIGN KEY (card_id) REFERENCES credit_cards (ID) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS credit_card_expense (
    ID uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    category_id int NOT NULL,
    amount double precision NOT NULL,
    description character varying(255),
    date date NOT NULL,
    card_id uuid NOT NULL,
    installment_amount double precision NOT NULL,
    installments_number int NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users (ID) ON DELETE CASCADE,
    CONSTRAINT fk_category_id FOREIGN KEY (category_id) REFERENCES categories (ID) ON DELETE RESTRICT,
    CONSTRAINT fk_card_id FOREIGN KEY (card_id) REFERENCES credit_cards (ID) ON DELETE SET NULL
);

-- Índices (exemplo)
CREATE INDEX idx_simple_expense_user_id ON simple_expense (user_id);
CREATE INDEX idx_recurring_expense_user_id ON recurring_expense (user_id);
CREATE INDEX idx_credit_card_expense_user_id ON credit_card_expense (user_id);


---- create above / drop below ----

DROP TABLE IF EXISTS expense;
DROP TABLE IF EXISTS recurring_expense;
DROP TABLE IF EXISTS credit_card_expense;