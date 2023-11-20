package models

type Sprint struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Judge     int    `json:"judge"`
	Completed bool   `json:"completed"`
	ProjectID int    `json:"projectId"`
}

type SprintMember struct {
	UserId      int    `json:"userId"`
	Designation string `json:"designation"`
}
