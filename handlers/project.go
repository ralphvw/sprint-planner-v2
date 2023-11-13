package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/ralphvw/sprint-planner-v2/helpers"
	"github.com/ralphvw/sprint-planner-v2/middleware"
	"github.com/ralphvw/sprint-planner-v2/services"
)

func AddProject(db *sql.DB) http.HandlerFunc {
	return middleware.CheckToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			helpers.HandleOptions(w, r)
			return
		}

		helpers.EnableCors(w)

		claims, ok := r.Context().Value("userClaims").(map[string]interface{})
		if !ok {
			helpers.LogAction("Error extracting user claims")
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		if r.Method == "GET" {
			helpers.LogAction(fmt.Sprintf("Claims part 2: %v", claims))
			helpers.LogAction(reflect.TypeOf(claims["id"]).String())
			ownerId, ok := claims["id"].(float64)
			if !ok {
				helpers.LogAction("Wrong type assertion")
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}

			results, err := services.FetchProjects(db, int(ownerId))
			if err != nil {
				helpers.LogAction("Get Projects: " + err.Error())
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}

			message := "Projects fetched successfully"
			helpers.LogAction("Projects fetched for userID: " + fmt.Sprintf("%d", int(ownerId)))

			helpers.SendResponse(w, r, message, results)
			return
		}

		var requestPayload struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			UserIDs     []int  `json:"userIds"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestPayload)
		if err != nil {
			helpers.LogAction("Bad Request: Add Project Handler" + err.Error())
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		requiredFields := []string{"Name", "Description", "UserIDs"}
		fieldsExist, missingFields := helpers.CheckFields(requestPayload, requiredFields)

		if !fieldsExist {
			helpers.LogAction(fmt.Sprintf("Missing fields: %v\n", missingFields))
			http.Error(w, fmt.Sprintf("Missing fields: %v\n", missingFields), http.StatusBadRequest)
			return
		}

		ownerId, ok := claims["id"].(float64)

		if !ok {
			helpers.LogAction("Wrong Type in Add Project Handler")
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		result, err := services.AddProject(db, requestPayload.Name, requestPayload.Description, int(ownerId))
		if err != nil {
			helpers.LogAction("Add Project Handler: " + err.Error())
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		for _, userId := range requestPayload.UserIDs {
			err := services.AddMember(db, result["id"], userId)
			if err != nil {
				helpers.LogAction("Add Project Member: " + err.Error())
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
		}

		message := "Project added successfully"

		helpers.LogAction("Project added successfully: " + "name: " + result["name"] + "id: " + result["id"])

		helpers.SendResponse(w, r, message, result)
	}))
}
