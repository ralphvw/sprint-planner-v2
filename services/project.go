package services

import (
	"database/sql"

	"github.com/ralphvw/sprint-planner-v2/models"
	"github.com/ralphvw/sprint-planner-v2/queries"
)

func AddProject(db *sql.DB, name string, description string, ownerId int) (map[string]string, error) {
	var project models.Project
	err := db.QueryRow(queries.AddProject, name, description, ownerId).Scan(&project.ID, &project.Name, &project.Description)
	if err != nil {
		return nil, err
	}
	result := map[string]string{
		"name":        name,
		"description": description,
	}
	return result, nil
}

func AddMember(db *sql.DB, projectId string, userId int) error {
	_, err := db.Exec(queries.AddMember, projectId, userId)
	return err
}

func FetchProjects(db *sql.DB, ownerId int) ([]models.Project, error) {
	var projects []models.Project

	rows, err := db.Query(queries.FetchProjects, ownerId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var project models.Project
		rows.Scan(&project.ID, &project.Name, &project.CreatedAt)
		projects = append(projects, project)
	}

	return projects, nil
}
