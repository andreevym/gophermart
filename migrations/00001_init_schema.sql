CREATE SEQUENCE IF NOT EXISTS users_id_seq;

CREATE TABLE IF NOT EXISTS users
(
    id         BIGINT PRIMARY KEY DEFAULT nextval('users_id_seq'),
    username   VARCHAR(255),
    password   VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255)       DEFAULT current_user
);

INSERT INTO users (username, password)
VALUES ('WithdrawUserID', 'qwe123');
INSERT INTO users (username, password)
VALUES ('AccrualUserID', 'qwe123');

CREATE TABLE IF NOT EXISTS orders
(
    number      VARCHAR(50) PRIMARY KEY,
    user_id     BIGINT,
    status      VARCHAR(255) NOT NULL,
    accrual     real,
    uploaded_at TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at  TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(255) NOT NULL DEFAULT current_user,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE SEQUENCE IF NOT EXISTS transactions_id_seq;

CREATE TABLE IF NOT EXISTS transactions
(
    transaction_id BIGINT PRIMARY KEY       DEFAULT nextval('transactions_id_seq'),
    from_user_id   BIGINT      NOT NULL,
    to_user_id     BIGINT      NOT NULL,
    amount         real        NOT NULL     DEFAULT 0 CHECK (amount >= 0),
    order_number   TEXT,
    created_at     TIMESTAMP WITH TIME ZONE                DEFAULT CURRENT_TIMESTAMP,
    created_by     VARCHAR(255)             DEFAULT current_user,
    operation_type VARCHAR(50) NOT NULL,
    date           TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_user_id) REFERENCES users (id),
    FOREIGN KEY (to_user_id) REFERENCES users (id)
);
