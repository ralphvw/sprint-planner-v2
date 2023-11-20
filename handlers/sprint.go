package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ralphvw/sprint-planner-v2/helpers"
	"github.com/ralphvw/sprint-planner-v2/middleware"
	"github.com/ralphvw/sprint-planner-v2/models"
	"github.com/ralphvw/sprint-planner-v2/queries"
	"github.com/ralphvw/sprint-planner-v2/services"
)

func AddSprint(db *sql.DB) http.HandlerFunc {
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
		userId, ok := claims["id"].(float64)
		if !ok {
			helpers.LogAction("Wrong type assertion for claims")
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		if r.Method == "GET" {
			projectId, err := strconv.Atoi(r.URL.Query().Get("projectId"))
			if err != nil {
				helpers.LogAction("Error converting query to int: " + err.Error())
				http.Error(w, "ProjectId missing", http.StatusBadRequest)
				return
			}
			page, err := strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				helpers.LogAction("Error converting query to int: " + err.Error())
				http.Error(w, "Page number missing", http.StatusBadRequest)
				return
			}

			e := services.CheckProjectOwner(db, int(userId), projectId)

			if e != nil {
				helpers.LogAction(fmt.Sprintf("User: %d trying to fetch sprints from project: %d  without membership", int(userId), projectId))
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			args := []interface{}{
				projectId,
			}

			var id int
			var name string

			keys := []string{
				"id",
				"name",
			}

			message := "Sprints fetched successfully"
			helpers.GetDataHandler(w, r, db, 10, page, queries.GetSprints, queries.CountSprints, message, args, keys, &id, &name)

		}

		var requestPayload struct {
			ProjectId int                   `json:"projectId"`
			Judge     int                   `json:"judge"`
			Name      string                `json:"name"`
			Members   []models.SprintMember `json:"members"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestPayload)
		if err != nil {
			helpers.LogAction("Bad Request: Add Project Handler" + err.Error())
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		requiredFields := []string{"ProjectId", "Name", "Judge"}
		fieldsExist, missingFields := helpers.CheckFields(requestPayload, requiredFields)

		if !fieldsExist {
			helpers.LogAction(fmt.Sprintf("Missing fields: %v\n", missingFields))
			http.Error(w, fmt.Sprintf("Missing fields: %v\n", missingFields), http.StatusBadRequest)
			return
		}

		e := services.CheckProjectOwner(db, int(userId), requestPayload.ProjectId)

		if e != nil {
			helpers.LogAction(fmt.Sprintf("Attempt to create sprint without ownership %d", int(userId)))
			http.Error(w, "Forbidden Request", http.StatusForbidden)
			return
		}

		result, err := services.AddSprint(db, requestPayload.Name, requestPayload.Judge, requestPayload.ProjectId)
		if err != nil {
			helpers.LogAction("Error adding sprint: " + err.Error())
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		for _, member := range requestPayload.Members {
			err := services.AddSprintMember(db, member.UserId, member.Designation, result["id"])
			if err != nil {
				helpers.LogAction("Error adding member to sprint " + err.Error())
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
		}

		message := "Sprint created successfully"
		helpers.LogAction(fmt.Sprintf("Sprint created successfully by: %d", int(userId)))

		helpers.SendResponse(w, r, message, result)
	}))
}
