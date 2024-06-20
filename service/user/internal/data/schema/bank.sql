CREATE DATABASE bank;

CREATE TABLE accounts
(
    id         bigserial primary key,
    owner      varchar     not null,
    balance    bigint      not null,
    currency   varchar     not null,
    created_at timestamptz not null default now()
);

CREATE TABLE entries
(
    id         bigserial primary key,
    account_id bigint REFERENCES accounts (id) not null,
    amount     bigint                          not null,
    created_at timestamptz                     not null default now()
);

CREATE TABLE transfers
(
    id              bigserial primary key,
    from_account_id bigint REFERENCES accounts (id) not null,
    to_account_id   bigint REFERENCES accounts (id) not null,
    amount          bigint                          not null,
    created_at      timestamptz                     not null default now()
);
