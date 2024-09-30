-- name: CreateAccount :one
INSERT INTO accounts (user_id, account_number, balance, currency)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, account_number, balance, currency;