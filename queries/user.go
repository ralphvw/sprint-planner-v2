package queries

var CreateUser string = `INSERT INTO users (first_name, last_name, email, hash) VALUES ($1, $2, $3, $4) RETURNING id, first_name, last_name, email`

var GetUserByEmail string = "SELECT id, email, hash, first_name, last_name FROM users WHERE email=$1"

var ResetPassword string = "UPDATE users SET hash=$1 WHERE email=$2 RETURNING email"

var SearchUsers string = "SELECT id, email, first_name, last_name FROM users WHERE email ILIKE $1 OR CONCAT(first_name, ' ', last_name) ILIKE $1"
