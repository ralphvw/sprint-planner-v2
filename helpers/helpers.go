package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ralphvw/sprint-planner-v2/models"
)

const (
	SERVER_ERROR = "Server error"
)

func LogAction(message string) {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log page")
	}

	defer file.Close()

	customFormat := "Mon 02 Jan, 2006 @ 15:04"
	logger := log.New(file, time.Now().Format(customFormat)+" ", 0)
	logger.Printf(message)
	fmt.Println(message)
}

func GetDataHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, pageSize int, page int, query string, countQuery string, message string, args []interface{}, destinations ...interface{}) {
	offset := (page - 1) * pageSize

	queryString := fmt.Sprintf("%s LIMIT %d OFFSET %d", query, pageSize, offset)

	rows, err := db.Query(queryString, args...)

	if err != nil {
		LogAction(err.Error())
		http.Error(w, SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		if err := rows.Scan(destinations...); err != nil {
			LogAction(err.Error())
			http.Error(w, SERVER_ERROR, http.StatusInternalServerError)
			return
		}

		result := make(map[string]interface{})

		for i := 0; i < len(destinations); i += 2 {
			key, ok := destinations[i].(string)
			if !ok {
				LogAction("Invalid key type")

				http.Error(w, SERVER_ERROR, http.StatusInternalServerError)
				return
			}
			value := destinations[i+1]
			result[key] = value
		}

		results = append(results, result)
	}
	totalRows := 0
	err = db.QueryRow(countQuery).Scan(&totalRows)
	if err != nil {
		LogAction(err.Error())
		http.Error(w, SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	totalPages := (totalRows + pageSize - 1) / pageSize

	response := models.Response{
		Message:    message,
		Data:       results,
		TotalPages: totalPages,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		LogAction(err.Error())
		http.Error(w, SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func GetSingleDataHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, query string, message string, args []interface{}, destinations ...interface{}) {
	row := db.QueryRow(query, args...)

	if err := row.Scan(destinations...); err != nil {
		LogAction(err.Error())
		http.Error(w, SERVER_ERROR, http.StatusInternalServerError)
		return
	}
	result := make(map[string]interface{})

	for i := 0; i < len(destinations); i += 2 {
		key, ok := destinations[i].(string)
		if !ok {
			LogAction("Invalid key type")
			http.Error(w, SERVER_ERROR, http.StatusInternalServerError)
			return
		}
		value := destinations[i+1]
		result[key] = value
	}
	response := models.Response{
		Message: message,
		Data:    result,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		LogAction(err.Error())
		http.Error(w, SERVER_ERROR, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)

}
