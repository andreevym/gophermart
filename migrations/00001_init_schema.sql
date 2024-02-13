CREATE SEQUENCE IF NOT EXISTS users_id_seq;

CREATE TABLE IF NOT EXISTS users
(
    id         BIGINT PRIMARY KEY DEFAULT nextval('users_id_seq'),
    username   VARCHAR(255),
    password   VARCHAR(255),
    created_at TIMESTAMP          DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255)       DEFAULT current_user
);

INSERT INTO users (username, password)
VALUES ('SystemUserID', 'qwe123');

CREATE SEQUENCE IF NOT EXISTS orders_id_seq;

CREATE TABLE IF NOT EXISTS orders
(
    id          BIGINT PRIMARY KEY DEFAULT nextval('orders_id_seq'),
    number      VARCHAR(255) UNIQUE,
    user_id     BIGINT,
    status      VARCHAR(255) NOT NULL,
    accrual     INT,
    uploaded_at TIMESTAMP,
    created_at  TIMESTAMP          DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255)       DEFAULT current_user,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS user_accounts
(
    user_id    BIGINT PRIMARY KEY,
    balance    BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255)    DEFAULT current_user,
    updated_at TIMESTAMP       DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255)    DEFAULT current_user,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE SEQUENCE IF NOT EXISTS transactions_id_seq;

CREATE TABLE IF NOT EXISTS transactions
(
    transaction_id BIGINT PRIMARY KEY       DEFAULT nextval('transactions_id_seq'),
    from_user_id   BIGINT      NOT NULL,
    to_user_id     BIGINT      NOT NULL,
    amount         BIGINT      NOT NULL     DEFAULT 0 CHECK (amount >= 0),
    reason         TEXT        NOT NULL,
    created_at     TIMESTAMP                DEFAULT CURRENT_TIMESTAMP,
    created_by     VARCHAR(255)             DEFAULT current_user,
    operation_type VARCHAR(50) NOT NULL,
    date           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user_id) REFERENCES users (id),
    FOREIGN KEY (to_user_id) REFERENCES users (id)
);
