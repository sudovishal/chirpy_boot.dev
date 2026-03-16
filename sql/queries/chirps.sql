-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES(gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: GetAllChirps :many
SELECT * from chirps
order by created_at;

-- name: GetChirpById :one
SELECT * from chirps where id= $1 order by created_at;

-- name: GetChirpByAuthor :many
select * from chirps where user_id = $1 order by created_at;

-- name: DeleteChirpByID :one
delete from chirps where user_id=$1 and id=$2
RETURNING *;
