package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ralphvw/sprint-planner-v2/helpers"
	"github.com/ralphvw/sprint-planner-v2/models"
	"github.com/ralphvw/sprint-planner-v2/queries"
)

func SearchUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			helpers.HandleOptions(w, r)
			return
		}

		helpers.EnableCors(w)
		searchTerm := r.URL.Query().Get("search")
		if searchTerm == "" {
			helpers.LogAction("Missing query param for Get all users")
			http.Error(w, "Missing 'search' query param", http.StatusBadRequest)
			return
		}
		var users []models.User

		rows, err := db.Query(queries.SearchUsers, "%"+searchTerm+"%")
		if err != nil {
			helpers.LogAction("Query failed: " + err.Error())
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var user models.User
			err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName)
			if err != nil {
				helpers.LogAction("Search Users Handler: " + err.Error())
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}

			users = append(users, user)
		}

		results := map[string]interface{}{
			"users": users,
		}

		message := "Users fetched successfully"

		helpers.LogAction("Users fetched")
		helpers.SendResponse(w, r, message, results)
	}
}
