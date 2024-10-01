-- name: CreateAccount :one
INSERT INTO accounts (user_id, account_number, balance, currency)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, account_number, balance, currency;

-- name: UpdateAccountBalance :exec
UPDATE accounts
SET balance = balance + $1
WHERE id = $2;

-- name : GetAccountByUserId :one
SELECT * FROM accounts WHERE user_id = $1;