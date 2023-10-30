package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ralphvw/sprint-planner-v2/helpers"
	"github.com/ralphvw/sprint-planner-v2/models"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			helpers.LogAction("Error: Login Handler: " + err.Error())
			http.Error(w, "Invalid Input", http.StatusBadRequest)
			return
		}

		authenticatedUser, err := helpers.AuthenticateUser(db, user.Email, user.Password)
		if err != nil {
			helpers.LogAction("Invalid Credentials" + user.Email)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := helpers.CreateToken(authenticatedUser.ID, authenticatedUser.FirstName, authenticatedUser.LastName, authenticatedUser.Email)
		if err != nil {
			helpers.LogAction("Token creation failed")
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		userResponse := models.UserResponse{
			ID:        authenticatedUser.ID,
			Email:     authenticatedUser.Email,
			FirstName: authenticatedUser.FirstName,
			LastName:  authenticatedUser.LastName,
		}

		result := make(map[string]interface{})
		result["token"] = token
		result["user"] = userResponse
		message := "Login Successful"
		response := models.SingleResponse{
			Message: message,
			Data:    result,
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			helpers.LogAction(err.Error())
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)

	}
}

func SignUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			helpers.LogAction("Error: Signup Handler: " + err.Error())
			http.Error(w, "Invalid Input", http.StatusBadRequest)
			return
		}

		plainTextPassword := user.Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), bcrypt.DefaultCost)

		if err != nil {
			helpers.LogAction("Error: Hashing Password " + err.Error())
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		user.Hash = hashedPassword
		query := `INSERT INTO users (first_name, last_name, email, hash) VALUES ($1, $2, $3, $4) RETURNING id, first_name as "firstName", last_name as "lastName", email as "email"`

		var newUser models.User

		err = db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Hash).Scan(&newUser.ID, &newUser.FirstName, &newUser.LastName, &newUser.Email)

		if err != nil {
			helpers.LogAction("Error: Failed to create user " + err.Error())
			http.Error(w, "User already exists", http.StatusInternalServerError)
			return
		}

		userResponse := models.UserResponse{
			ID:        newUser.ID,
			FirstName: newUser.FirstName,
			LastName:  newUser.LastName,
			Email:     newUser.Email,
		}

		response := models.SingleResponse{
			Message: "User created successfully",
			Data:    userResponse,
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			helpers.LogAction(err.Error())
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		helpers.LogAction("User created successfully: " + newUser.Email)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(responseJSON)

	}
}
