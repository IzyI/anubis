package DTO

type CreateProjectValid struct {
	ProjectName string `json:"project_name"  binding:"required,min=12,max=64"`
}

type UpdateProjectValid struct {
	ProjectName string `json:"project_name"  binding:"required,min=12,max=64"`
	ID          string `uri:"id"  binding:"required,min=23,max=25"`
}

type DelProjectValid struct {
	ID string `uri:"id"  binding:"required,min=23,max=25"`
}

type AnswerProject struct {
	ProjectName string `json:"project_name"`
	ProjectId   string `json:"project_id"`
}
