-- +goose Up

-- Create wallets table
CREATE TABLE users (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    is_identified BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE wallets (
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id uuid NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES users(id)
);

-- Create transactions table
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    wallet_id uuid NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_wallet_id FOREIGN KEY(wallet_id) REFERENCES wallets(id)
);

-- Insert sample data
INSERT INTO users (is_identified) VALUES
(FALSE),
(TRUE),
(FALSE),
(TRUE);

INSERT INTO wallets (user_id, balance) VALUES
('7a9abd33-7b21-482c-8c3b-d256e711170e', 500000),
('d3211ea3-fdc1-4625-8f1d-a6fa40f9f4ce', 300000),
('a29cfe7e-d2dd-4c96-826c-c94d96569fe5', 800000),
('3ee68254-1f2a-4d9a-b6f5-09d27ee43445', 0);

INSERT INTO transactions (wallet_id, amount, created_at) VALUES
('9e323a19-73bf-4600-a049-7397dcea5751', 100000, CURRENT_TIMESTAMP - INTERVAL '5 days'),
('9e323a19-73bf-4600-a049-7397dcea5751', 200000, CURRENT_TIMESTAMP - INTERVAL '3 days'),
('4aa3e422-eff1-46de-adab-8ea0c3893fe7', 150000, CURRENT_TIMESTAMP - INTERVAL '2 days'),
('4f09c559-5ab1-427f-bb51-8455bae3a067', 300000, CURRENT_TIMESTAMP - INTERVAL '1 day');

-- +goose Down
drop table transactions;
drop table wallets;
drop table users;