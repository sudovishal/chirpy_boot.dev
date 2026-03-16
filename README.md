# chirpy_boot.dev

# Goose Migration Up
goose postgres <connection_string> up

example:
`goose postgres "postgres://vishals:pgsql@localhost:5432/chirpy" up`

# Goose Migration Down
goose postgres <connection_string> down

goose postgres "postgres://vishals:pgsql@localhost:5432/chirpy" down-to 20260123042621


example:
`goose postgres "postgres://vishals:pgsql@localhost:5432/chirpy" down`

# SQLC Generate

`sqlc generate`

    
export PATH="/usr/local/opt/postgresql@15/bin:$PATH"



pgcli -h localhost -p 5432 -U vishals -d chirpy


I am getting an error when I hit a /api/chirps POST endpoint with the token which is generated after hitting the /api/refresh endpoint.
