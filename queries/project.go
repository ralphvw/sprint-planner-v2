package queries

var AddProject string = "INSERT INTO projects(name, description, owner_id) VALUES($1, $2, $3) RETURNING id, name, description"
var AddMember string = "INSERT INTO project_members(project_id, user_id) VALUES($1, $2)"

var FetchProjects string = "SELECT id, name, created_at from projects WHERE owner_id=$1"

var CountProjects string = "SELECT COUNT(*) FROM projects WHERE owner_id=$1"

var CheckProjectMember string = "SELECT id FROM project_members WHERE user_id=$1 AND project_id=$2"

var CheckProjectOwner string = "SELECT owner_id FROM projects WHERE id=$1"

var GetSingleProject string = `SELECT p.id, p.name, p.description, p.created_at,
(
SELECT COUNT(*) FROM project_members
WHERE project_id=$1
) as "memberCount",
CASE
  WHEN COUNT(u.id) = 0 THEN NULL
  ELSE
  json_agg(
  json_build_object(
    'id', u.id,
    'firstName', u.first_name,
    'lastName', u.last_name
  )
)  END AS members,
json_build_object(
  'id', o.id,
  'firstName', o.first_name,
  'lastName', o.last_name
) as owner
FROM projects p
 LEFT JOIN project_members pm on pm.project_id = p.id
 LEFT JOIN users u on u.id = pm.user_id
 LEFT JOIN users o on o.id = p.owner_id
 WHERE p.id=$1
 GROUP BY p.id, p.name, p.description, p.created_at, o.id, o.first_name, o.last_name
`
