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

			search := r.URL.Query().Get("search")
			searchTerm := "%" + search + "%"

			args := []interface{}{
				int(ownerId),
				searchTerm,
			}

			var id int
			var name string
			var createdAt time.Time
			var plans int

			keys := []string{
				"id",
				"name",
				"createdAt",
				"plans",
			}

			message := "Projects fetched successfully"
			helpers.GetDataHandler(w, r, db, 10, page, queries.FetchProjects, queries.CountProjects, message, args, keys, &id, &name, &createdAt, &plans)
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

		res, ok := result["id"].(int)

		if !ok {
			helpers.LogAction("Error converting projectId into an int")
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		for _, userId := range requestPayload.UserIDs {
			err := services.AddMember(db, res, userId)
			if err != nil {
				helpers.LogAction("Add Project Member: " + err.Error())
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}
		}

		message := "Project added successfully"

		helpers.LogAction("Project added successfully")

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
				http.Error(w, "Forbidden Request", http.StatusForbidden)
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

func AddMember(db *sql.DB) http.HandlerFunc {
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

		var requestPayload struct {
			ProjectId int `json:"projectId"`
			UserId    int `json:"userId"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestPayload)
		if err != nil {
			helpers.LogAction("Bad Request: Add Project Handler" + err.Error())
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		requiredFields := []string{"ProjectId", "UserId"}
		fieldsExist, missingFields := helpers.CheckFields(requestPayload, requiredFields)

		if !fieldsExist {
			helpers.LogAction(fmt.Sprintf("Missing fields: %v\n", missingFields))
			http.Error(w, fmt.Sprintf("Missing fields: %v\n", missingFields), http.StatusBadRequest)
			return
		}

		e := services.CheckProjectOwner(db, int(userId), requestPayload.ProjectId)

		if e != nil {
			helpers.LogAction(fmt.Sprintf("Attempt to view project without membership %d", int(userId)))
			http.Error(w, "Forbidden Request", http.StatusForbidden)
			return
		}

		if r.Method == "DELETE" {
			er := services.DeleteMember(db, requestPayload.UserId, requestPayload.ProjectId)
			if er != nil {
				helpers.LogAction("Error deleting member " + err.Error())
				http.Error(w, "Server Error", http.StatusInternalServerError)
				return
			}

			result := map[string]interface{}{
				"userId":    requestPayload.UserId,
				"projectId": requestPayload.ProjectId,
			}

			message := "Member removed sucessfully"
			helpers.LogAction("Member removed successfully to project")
			helpers.SendResponse(w, r, message, result)
			return
		}

		errr := services.CheckProjectMember(db, requestPayload.UserId, requestPayload.ProjectId)

		if errr == nil {
			helpers.LogAction("Attempt to add an already existing member")
			http.Error(w, "User is already a project member", http.StatusConflict)
			return
		}
		er := services.AddMember(db, requestPayload.ProjectId, requestPayload.UserId)
		if er != nil {
			helpers.LogAction("Add Project Member: " + err.Error())
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		result := map[string]interface{}{
			"userId":    requestPayload.UserId,
			"projectId": requestPayload.ProjectId,
		}

		message := "Member added sucessfully"
		helpers.LogAction("Member added successfully to project")
		helpers.SendResponse(w, r, message, result)

	}))
}

func GetMembers(db *sql.DB) http.HandlerFunc {
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

		projectId, err := strconv.Atoi(r.URL.Query().Get("projectId"))
		if err != nil {
			helpers.LogAction("Error converting query to int: " + err.Error())
			http.Error(w, "Page number missing", http.StatusBadRequest)
			return
		}
		e := services.CheckProjectMember(db, int(userId), projectId)

		if e != nil {
			helpers.LogAction(fmt.Sprintf("Attempt to view project without membership %d", int(userId)))
			http.Error(w, "Forbidden Request", http.StatusForbidden)
			return
		}

		result, err := services.GetProjectMembers(db, projectId)

		if err != nil {
			helpers.LogAction("Get Project Members Handler: " + err.Error())
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}

		message := "Members fetched successfully"

		helpers.LogAction(fmt.Sprintf("Members fetched for: %d", projectId))
		helpers.SendResponse(w, r, message, result)

	}))

}
