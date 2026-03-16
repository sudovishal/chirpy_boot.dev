-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email,hashed_password)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: DeleteAllUsers :exec
TRUNCATE TABLE users, chirps CASCADE;

-- name: GetUserByEmail :one
SELECT * from users where email=$1;

-- name: GenerateRefreshToken :one
INSERT INTO refresh_tokens (token,created_at,updated_at,user_id,expires_at,revoked_at)
VALUES($1,NOW(), NOW(), $2, $3, NULL)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM users
JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND refresh_tokens.revoked_at IS NULL
AND refresh_tokens.expires_at > NOW();

-- name: DeleteRFToken :exec
UPDATE refresh_tokens SET revoked_at = NOW(), updated_at = NOW() where token=$1;

-- name: UpdateEmailPass :one
UPDATE users SET email=$1,hashed_password=$2, updated_at = NOW() where id=$3
RETURNING *;

-- name: UpdateRedMembership :one
UPDATE users SET is_chirpy_red = true WHERE id = $1
RETURNING *;
