package services

import (
	"database/sql"
	"errors"

	"github.com/ralphvw/sprint-planner-v2/models"
	"github.com/ralphvw/sprint-planner-v2/queries"
)

func AddProject(db *sql.DB, name string, description string, ownerId int) (map[string]interface{}, error) {
	var project models.Project
	err := db.QueryRow(queries.AddProject, name, description, ownerId).Scan(&project.ID, &project.Name, &project.Description)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{
		"id":          project.ID,
		"name":        name,
		"description": description,
	}
	return result, nil
}

func AddMember(db *sql.DB, projectId int, userId int) error {
	_, err := db.Exec(queries.AddMember, projectId, userId)
	return err
}

func DeleteMember(db *sql.DB, userId int, projectId int) error {
	_, err := db.Exec(queries.RemoveMember, userId, projectId)
	return err
}

func CheckProjectMember(db *sql.DB, userId int, projectId int) error {
	var ownerId int
	var id int

	db.QueryRow(queries.CheckProjectOwner, projectId).Scan(&ownerId)

	err := db.QueryRow(queries.CheckProjectMember, userId, projectId).Scan(&id)

	if err == sql.ErrNoRows && userId != ownerId {
		return err
	}

	return nil
}

func CheckProjectOwner(db *sql.DB, userId int, projectId int) error {
	var ownerId int

	db.QueryRow(queries.CheckProjectOwner, projectId).Scan(&ownerId)

	if userId != ownerId {
		return errors.New("Forbidden")
	}

	return nil
}
