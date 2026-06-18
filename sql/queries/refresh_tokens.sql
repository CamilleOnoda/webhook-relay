-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, user_id, expires_at)
VALUES($1, $2, $3)
RETURNING *;

-- name: RevokeRefreshToken :execrows
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1
    AND revoked_at IS NULL;
