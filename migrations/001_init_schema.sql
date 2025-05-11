-- Создание расширения для шифрования
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Таблица пользователей
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(100) NOT NULL UNIQUE,
                       email VARCHAR(255) NOT NULL UNIQUE,
                       password_hash VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица банковских счетов
CREATE TABLE accounts (
                          id SERIAL PRIMARY KEY,
                          user_id INTEGER NOT NULL REFERENCES users(id),
                          number VARCHAR(20) NOT NULL UNIQUE,
                          balance NUMERIC(15, 2) NOT NULL DEFAULT 0.0,
                          type VARCHAR(20) NOT NULL, -- 'debit', 'credit', 'saving'
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица банковских карт
CREATE TABLE cards (
                       id SERIAL PRIMARY KEY,
                       account_id INTEGER NOT NULL REFERENCES accounts(id),
                       number_encrypted BYTEA NOT NULL,
                       expiry_month_encrypted BYTEA NOT NULL,
                       expiry_year_encrypted BYTEA NOT NULL,
                       cvv_hash VARCHAR(255) NOT NULL,
                       number_hmac VARCHAR(255) NOT NULL,
                       expiry_hmac VARCHAR(255) NOT NULL,
                       is_active BOOLEAN NOT NULL DEFAULT TRUE,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Таблица транзакций
CREATE TABLE transactions (
                              id SERIAL PRIMARY KEY,
                              from_account_id INTEGER REFERENCES accounts(id),
                              to_account_id INTEGER REFERENCES accounts(id),
                              amount NUMERIC(15, 2) NOT NULL,
                              type VARCHAR(20) NOT NULL, -- 'deposit', 'withdrawal', 'transfer', 'payment'
                              description TEXT,
                              status VARCHAR(20) NOT NULL DEFAULT 'completed', -- 'pending', 'completed', 'failed'
                              created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              CHECK (amount > 0),
                              CHECK ((from_account_id IS NOT NULL) OR (to_account_id IS NOT NULL))
);

-- Таблица кредитов
CREATE TABLE credits (
                         id SERIAL PRIMARY KEY,
                         user_id INTEGER NOT NULL REFERENCES users(id),
                         account_id INTEGER NOT NULL REFERENCES accounts(id),
                         amount NUMERIC(15, 2) NOT NULL,
                         interest_rate NUMERIC(5, 2) NOT NULL,
                         term_months INTEGER NOT NULL,
                         monthly_payment NUMERIC(15, 2) NOT NULL,
                         remaining_amount NUMERIC(15, 2) NOT NULL,
                         status VARCHAR(20) NOT NULL DEFAULT 'active', -- 'active', 'closed', 'defaulted'
                         created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                         updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                         CHECK (amount > 0),
                         CHECK (interest_rate > 0),
                         CHECK (term_months > 0),
                         CHECK (monthly_payment > 0)
);

-- Таблица графика платежей
CREATE TABLE payment_schedules (
                                   id SERIAL PRIMARY KEY,
                                   credit_id INTEGER NOT NULL REFERENCES credits(id),
                                   payment_date DATE NOT NULL,
                                   payment_amount NUMERIC(15, 2) NOT NULL,
                                   principal_amount NUMERIC(15, 2) NOT NULL,
                                   interest_amount NUMERIC(15, 2) NOT NULL,
                                   is_paid BOOLEAN NOT NULL DEFAULT FALSE,
                                   created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                   updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                   CHECK (payment_amount > 0),
                                   CHECK (principal_amount > 0),
                                   CHECK (interest_amount > 0)
);

-- Индексы для ускорения запросов
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_cards_account_id ON cards(account_id);
CREATE INDEX idx_transactions_from_account_id ON transactions(from_account_id);
CREATE INDEX idx_transactions_to_account_id ON transactions(to_account_id);
CREATE INDEX idx_credits_user_id ON credits(user_id);
CREATE INDEX idx_credits_account_id ON credits(account_id);
CREATE INDEX idx_payment_schedules_credit_id ON payment_schedules(credit_id);
CREATE INDEX idx_payment_schedules_payment_date ON payment_schedules(payment_date);
