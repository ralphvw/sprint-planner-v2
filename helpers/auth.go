package helpers

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ralphvw/sprint-planner-v2/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func AuthenticateUser(db *sql.DB, email string, password string) (*models.User, error) {
	query := "SELECT id, email, hash FROM users WHERE email=$1"
	row := db.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Hash)
	if err != nil {
		if err == sql.ErrNoRows {
			LogAction("User not found" + email)
			return nil, fmt.Errorf("User not found")
		}
	}

	if !comparePasswords(password, string(user.Hash)) {
		LogAction("Authentication failed" + email)
		return nil, fmt.Errorf("Authentication Failed")
	}

	return &user, nil
}

func comparePasswords(plain, hash string) bool {
	hashBytes := []byte(hash)
	err := bcrypt.CompareHashAndPassword(hashBytes, []byte(plain))
	return err == nil
}

func CreateToken(userId int, firstName string, lastName string, email string) (string, error) {
	claims := &models.Claims{
		UserID:    userId,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
