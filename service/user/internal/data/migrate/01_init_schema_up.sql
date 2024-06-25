-- CREATE DATABASE bank;

-- 账户
CREATE TABLE accounts
(
    id         bigserial primary key,
    owner      varchar     not null,
    balance    bigint      not null,
    currency   varchar     not null,
    created_at timestamptz not null default now()
);

-- 账单
CREATE TABLE entries
(
    id         bigserial primary key,
    account_id bigint      not null,
    amount     bigint      not null,
    created_at timestamptz not null default now()
);

-- 转账
CREATE TABLE transfers
(
    id              bigserial primary key,
    from_account_id bigint      not null,
    to_account_id   bigint      not null,
    amount          bigint      not null,
    created_at      timestamptz not null default now()
);

-- 外键约束
ALTER TABLE entries
    ADD FOREIGN KEY (account_id) REFERENCES accounts (id);
ALTER TABLE transfers
    ADD FOREIGN KEY (from_account_id) REFERENCES accounts (id);
ALTER TABLE transfers
    ADD FOREIGN KEY (to_account_id) REFERENCES accounts (id);
