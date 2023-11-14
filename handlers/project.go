package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ralphvw/sprint-planner-v2/helpers"
	"github.com/ralphvw/sprint-planner-v2/middleware"
	"github.com/ralphvw/sprint-planner-v2/queries"
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
			page, err := strconv.Atoi(r.URL.Query().Get("page"))
			if err != nil {
				helpers.LogAction("Error converting query to int: " + err.Error())
				http.Error(w, "Page number missing", http.StatusBadRequest)
				return
			}
			ownerId, ok := claims["id"].(float64)
			if !ok {
				helpers.LogAction("Wrong type assertion")
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}

			args := []interface{}{
				int(ownerId),
			}

			var id int
			var name string
			var createdAt time.Time

			keys := []string{
				"id",
				"name",
				"createdAt",
			}

			message := "Projects fetched successfully"
			helpers.GetDataHandler(w, r, db, 10, page, queries.FetchProjects, queries.CountProjects, message, args, keys, &id, &name, &createdAt)
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

func SingleProject(db *sql.DB) http.HandlerFunc {
	return middleware.CheckToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			helpers.HandleOptions(w, r)
			return
		}

		claims, ok := r.Context().Value("userClaims").(map[string]interface{})

		if !ok {
			helpers.LogAction("Errors extracting user claims")
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		if r.Method == "GET" {
			userId, ok := claims["id"].(float64)
			if !ok {
				helpers.LogAction("Wrong type assertion for claims")
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}

			projectId := strings.TrimPrefix(r.URL.Path, "/project/")

			if projectId == "" {
				helpers.LogAction("Invalid URL Path: Missing id argument")
				http.Error(w, "Invalid URL Pah: Missing id argument", http.StatusBadRequest)
				return
			}

			projectID, err := strconv.Atoi(projectId)

			if err != nil {
				helpers.LogAction("Error typecasting to int")
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}

			e := services.CheckProjectMember(db, int(userId), projectID)

			if e != nil {
				helpers.LogAction(fmt.Sprintf("Attempt to view project without membership %d", int(userId)))
				http.Error(w, "Unauthorized Request", http.StatusUnauthorized)
				return
			}

			message := "Project fetched successfully"

			var id int
			var name string
			var description string
			var createdAt time.Time
			var memberCount int
			var members *json.RawMessage
			var owner json.RawMessage

			keys := []string{
				"id",
				"name",
				"description",
				"createdAt",
				"memberCount",
				"members",
				"owner",
			}

			args := []interface{}{
				projectID,
			}

			helpers.GetSingleDataHandler(w, r, db, queries.GetSingleProject, message, args, keys, &id, &name, &description, &createdAt, &memberCount, &members, &owner)
		}
	}))
}
