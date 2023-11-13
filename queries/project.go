package queries

var AddProject string = "INSERT INTO projects(name, description, owner_id) VALUES($1, $2, $3) RETURNING id, name, description"
var AddMember string = "INSERT INTO project_members(project_id, user_id) VALUES($1, $2)"

var FetchProjects string = "SELECT id, name, created_at from projects WHERE owner_id=$1"
