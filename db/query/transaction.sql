-- name: Deposit :one
INSERT INTO transactions (account_id, transaction_type, amount, status, reference_id)
VALUES ($1, 'deposit', $2, 'pending', $3)
RETURNING id, account_id, transaction_type, amount, status, reference_id;

