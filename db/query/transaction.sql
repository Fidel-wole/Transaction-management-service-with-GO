-- name: Deposit :one
INSERT INTO transactions (account_id, transaction_type, amount, status, reference_id)
VALUES ($1, 'deposit', $2, 'pending', $3)
RETURNING id, account_id, transaction_type, amount, status, reference_id;

-- name: UpdateSenderBalance :exec
UPDATE accounts
SET balance = balance - $2
WHERE account_number = $1;

-- name: UpdateReceiverBalance :exec
UPDATE accounts
SET balance = balance + $2
WHERE account_number = $1;

-- name: CreateDebitTransaction :exec
INSERT INTO transactions (account_id, transaction_type, amount, status, reference_id)
VALUES ($1, 'withdraw', $2, 'completed', $3);

-- name: CreateCreditTransaction :exec
INSERT INTO transactions (account_id, transaction_type, amount, status, reference_id)
VALUES ($1, 'deposit', $2, 'completed', $3);

-- name: GetAccountIDByAccountNumber :one
SELECT id
FROM accounts
WHERE account_number = $1;
