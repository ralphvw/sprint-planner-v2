package services

import (
	"database/sql"

	"github.com/ralphvw/sprint-planner-v2/models"
	"github.com/ralphvw/sprint-planner-v2/queries"
)

func AddSprint(db *sql.DB, name string, judge int, projectId int) (*map[string]interface{}, error) {
	var sprint models.Sprint

	err := db.QueryRow(queries.CreateSprint, name, judge, projectId).Scan(&sprint.ID, &sprint.Name, &sprint.Judge, &sprint.ProjectID)

	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"id":        sprint.ID,
		"name":      sprint.Name,
		"judge":     sprint.Judge,
		"projectId": sprint.ProjectID,
	}

	return &result, nil
}

func AddSprintMember(db *sql.DB, userId int, designation string, sprintId interface{}) error {
	_, err := db.Exec(queries.AddSprintMember, userId, designation, sprintId)
	if err != nil {
		return err
	}
	return nil
}
