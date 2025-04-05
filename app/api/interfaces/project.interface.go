package interfaces

import (
	"anubis/app/DTO"
	"anubis/app/core/common"
	"github.com/gin-gonic/gin"
)

type ProjectsUseCase interface {
	GetProjectsFlow(ctx *gin.Context, input *common.Empty) (DTO.AnswerProjectList, error)
	PostProjectsFlow(ctx *gin.Context, input *DTO.CreateProjectValid) (DTO.AnswerProject, error)
	UpdateProjectsFlow(ctx *gin.Context, uri *DTO.UriIDValid, body *DTO.UpdateProjectValid) (common.Empty, error)
	DelProjectsFlow(ctx *gin.Context, uri *DTO.UriIDValid) (common.Empty, error)
	//ChangeProjectOwner(ctx *gin.Context, input *DTO.DelProjectValid) (common.Empty, error)
}

type ProjectMembersUseCase interface {
	GetProjectMembersFlow(ctx *gin.Context, input *DTO.UriIDValid) (DTO.AnswerProjectID, error)
	PostProjectMembersFlow(ctx *gin.Context, uri *DTO.UriIDValid, body *DTO.MembersAddIdValid) (DTO.AnswerMembers, error)
	PutProjectMembersFlow(ctx *gin.Context, uri *DTO.UriIDValid, body *DTO.MembersAddIdValid) (DTO.AnswerMembers, error)
	DelProjectMembersFlow(ctx *gin.Context, uri *DTO.UriIDValid, body *DTO.MembersAddIdValid) (common.Empty, error)
}
