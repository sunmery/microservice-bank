CREATE DATABASE bank;

CREATE TABLE accounts
(
    id         bigserial PRIMARY KEY,
    owner      varchar     NOT NULL,
    balance    bigint      NOT NULL,
    currency   varchar     NOT NULL,
    created_at timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE entries
(
    id         bigserial PRIMARY KEY,
    account_id bigint REFERENCES accounts (id) NOT NULL,
    amount     bigint                          NOT NULL,
    created_at timestamptz                     NOT NULL DEFAULT (now())
);

CREATE TABLE transfers
(
    id              bigserial PRIMARY KEY,
    from_account_id bigint REFERENCES accounts (id) NOT NULL,
    to_account_id   bigint REFERENCES accounts (id) NOT NULL,
    amount          bigint                          NOT NULL,
    created_at      timestamptz                     NOT NULL DEFAULT (now())
);

CREATE INDEX idx_owner ON accounts (owner);
CREATE INDEX idx_account_id ON entries (account_id);
CREATE INDEX idx_from_account_id ON transfers (from_account_id);
CREATE INDEX idx_to_account_id ON transfers (to_account_id);
CREATE INDEX idx_from_to_account_id ON transfers (from_account_id, to_account_id);
