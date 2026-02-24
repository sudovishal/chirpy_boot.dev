-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email,hashed_password)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: DeleteAllUsers :exec
TRUNCATE TABLE users, chirps;

-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: GetAllChirps :many
SELECT * from chirps
order by created_at;

-- name: GetChirpById :one
SELECT * from chirps where id= $1;

-- name: GetUserByEmail :one
SELECT * from users where email=$1;


-- name: GenerateRefreshToken :one
INSERT INTO refresh_tokens (token,created_at,updated_at,user_id,expires_at,revoked_at)
VALUES($1,NOW(), NOW(), $2, $3, NULL)
RETURNING *;


-- name: GetRefreshToken :one
select * from refresh_tokens where token=$1;


-- name: GetUserFromRefreshToken :one
SELECT user_id from refresh_tokens where token=$1;
