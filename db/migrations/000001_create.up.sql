CREATE TABLE IF NOT EXISTS currencies (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name varchar(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS transfers (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    from_balance_id BIGINT UNSIGNED NOT NULL,
    to_balance_id BIGINT UNSIGNED NOT NULL,
    amount BIGINT NOT NULL COMMENT 'can be only positive'
);

CREATE TABLE IF NOT EXISTS balances (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT UNSIGNED NOT NULL,
    currency_id BIGINT UNSIGNED NOT NULL,
    amount BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS entries (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    balance_id BIGINT UNSIGNED NOT NULL,
    amount bigint NOT NULL COMMENT 'can be negative or positive'
);

CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    username varchar(255) NOT NULL
);

CREATE UNIQUE INDEX currencies_index_0 ON currencies(`name`);

CREATE UNIQUE INDEX users_index_1 ON users(`username`);

ALTER TABLE transfers ADD FOREIGN KEY (from_balance_id) REFERENCES balances(`id`);

ALTER TABLE transfers ADD FOREIGN KEY (to_balance_id) REFERENCES balances(`id`);

ALTER TABLE balances ADD FOREIGN KEY (user_id) REFERENCES users(`id`);

ALTER TABLE balances ADD FOREIGN KEY (currency_id) REFERENCES currencies(`id`);

ALTER TABLE entries ADD FOREIGN KEY (balance_id) REFERENCES balances(`id`);