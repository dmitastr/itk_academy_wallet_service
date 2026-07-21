CREATE TABLE IF NOT EXISTS wallet_transactions (
     id              serial primary key,
     wallet_id       UUID NOT NULL,
     amount          INT NOT NULL,
     created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_wallet_transactions_wallet_id ON wallet_transactions(wallet_id);