package services

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/ralphvw/sprint-planner-v2/helpers"
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

func GetProjectMembers(db *sql.DB, projectId int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	var user models.User

	rows, err := db.Query(queries.GetProjectMembers, projectId)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email)
		if err != nil {
			return nil, err
		}

		result := map[string]interface{}{
			"id":        reflect.Indirect(reflect.ValueOf(user.ID)).Interface(),
			"firstName": reflect.Indirect(reflect.ValueOf(user.FirstName)).Interface(),
			"lastName":  helpers.GetPointerValue(user.LastName),
		}

		helpers.LogAction(fmt.Sprintf("%v", result))

		results = append(results, result)
	}

	var owner models.User
	db.QueryRow(queries.GetProjectOwner, projectId).Scan(&owner.ID, &owner.FirstName, &owner.LastName, &owner.Email)

	ownerResult := map[string]interface{}{
		"id":        owner.ID,
		"firstName": owner.FirstName,
		"lastName":  owner.LastName,
	}

	results = append(results, ownerResult)
	return results, nil
}
