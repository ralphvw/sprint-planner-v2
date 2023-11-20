package queries

var CreateSprint string = `INSERT INTO sprints (name, judge, project_id) VALUES ($1, $2, $3) RETURNING id, name, judge, project_id`

var AddSprintMember string = `INSERT INTO sprint_members (user_id, designation, sprint_id) VALUES ($1, $2, $3)`

var GetSprints string = `SELECT id, name FROM sprints WHERE project_id=$1`

var CountSprints string = `SELECT COUNT(*) FROM sprints WHERE project_id=$1`
