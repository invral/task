BEGIN;

CREATE TABLE IF NOT EXISTS public.account (
    id SERIAL PRIMARY KEY NOT NULL,
    currency VARCHAR(3) NOT NULL,
    balance FLOAT NOT NULL DEFAULT 0,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS public.transaction (
       id SERIAL PRIMARY KEY NOT NULL,
       status VARCHAR(20) NOT NULL,
        account_id INT,
        amount INT NOT NULL DEFAULT 0,
       currency VARCHAR(3) NOT NULL,
        to_account INT
);

INSERT INTO account (id, currency, balance, password, email)
VALUES (1, 'USD', 100, 'qwerty1', '1@ya.ru');

INSERT INTO account (id, currency, balance, password, email)
VALUES (2, 'EUR', 200, 'qwerty2', '2@ya.ru');

INSERT INTO account (id, currency, balance, password, email)
VALUES (3, 'RUB', 1000, 'qwerty3', '3@ya.ru');


-- Транзакции, где to_account == 0 - это пополнение счета
INSERT INTO transaction (id, status, account_id, amount, currency, to_account)
VALUES (1, 'created', 1, 10, 'RUB', 0);

INSERT INTO transaction (id, status, account_id, amount, currency, to_account)
VALUES (2, 'created', 2, 10, 'USD', 0);

INSERT INTO transaction (id, status, account_id, amount, currency, to_account)
VALUES (3, 'created', 3, 5, 'USD', 1);


END;