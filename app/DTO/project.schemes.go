package DTO

type CreateProjectValid struct {
	ProjectName string `json:"project_name"  binding:"required,safe_text,min=12,max=64"`
}

type UpdateProjectValid struct {
	ProjectName string `json:"project_name" binding:"required,safe_text,min=12,max=64"`
}

type AnswerProject struct {
	ProjectName string `json:"project_name"`
	ProjectId   string `json:"project_id"`
}

type AnswerProjectList struct {
	ListProjects map[string]string `json:"list_projects"`
}

type UriIDValid struct {
	ID string `uri:"id"   binding:"required,object_id"`
}

type MembersAddIdValid struct {
	UserID string `json:"user_id"  binding:"required,object_id"`
	Role   string `json:"role"  binding:"omitempty,max=24"`
}

type AnswerMembers struct {
	UserID   string `json:"user_id"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}
type AnswerProjectID struct {
	ID      string          `json:"id"`
	Name    string          `json:"name"`
	Domain  string          `json:"domain"`
	Members []AnswerMembers `json:"string"`
}
