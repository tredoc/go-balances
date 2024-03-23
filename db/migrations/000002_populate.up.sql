INSERT INTO currencies(name)
    VALUES ('USD'), ('EUR'), ('UAH'), ('USDT');

INSERT INTO users(username)
    VALUES ('Kolya'), ('Vadym'), ('Serhii');

INSERT INTO balances(user_id, currency_id, amount)
    VALUES (1, 1, 100), (1, 2, 1000), (1, 3, 0), (1, 4, 300),
           (2, 1, 100), (2, 2, 1000), (2, 3, 0), (2, 4, 600),
           (3, 1, 100), (3, 2, 1000), (3, 3, 0), (3, 4, 900);

INSERT INTO entries (balance_id, amount)
    VALUES (1, 100), (2, 1000), (3, 0), (4, 600),
           (5, 100), (6, 1000), (7, 0), (8, 600),
           (9, 100), (10, 1000), (11, 0), (12, 600);

INSERT INTO transfers (from_balance_id, to_balance_id, amount)
    VALUES (4, 12, 300);
