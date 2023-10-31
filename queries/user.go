package queries

var CreateUser string = `INSERT INTO users (first_name, last_name, email, hash) VALUES ($1, $2, $3, $4) RETURNING id, first_name, last_name, email`

var GetUserByEmail string = "SELECT id, email, hash, first_name, last_name FROM users WHERE email=$1"
