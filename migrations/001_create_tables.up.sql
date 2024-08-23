-- Create wallets table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    is_identified BOOLEAN NOT NULL DEFAULT FALSE
)

CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER FOREIGN KEY (id) REFERENCES users,
    balance BIGINT NOT NULL DEFAULT 0,
);

-- Create transactions table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    wallet_id INTEGER REFERENCES wallets(id),
    amount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data
INSERT INTO users (id, is_identified) VALUES
(1, FALSE),
(2, TRUE),
(3, FALSE),
(4, TRUE);

INSERT INTO wallets (user_id, balance) VALUES
(1, 500000),
(2, 300000),
(3, 800000),
(4, 0);

INSERT INTO transactions (wallet_id, amount, timestamp) VALUES
(1, 100000, CURRENT_TIMESTAMP - INTERVAL '5 days'),
(1, 200000, CURRENT_TIMESTAMP - INTERVAL '3 days'),
(2, 150000, CURRENT_TIMESTAMP - INTERVAL '2 days'),
(3, 300000, CURRENT_TIMESTAMP - INTERVAL '1 day');