package interfaces

import (
	"anubis/app/api/DTO"
	"anubis/app/core/common"
	"github.com/gin-gonic/gin"
)

type ProjectsUseCase interface {
	GetProjectsFlow(ctx *gin.Context, input *common.Empty) (DTO.AnswerProjectList, error)
	PostProjectsFlow(ctx *gin.Context, input *DTO.CreateProjectValid) (DTO.AnswerProject, error)
	PutProjectsFlow(ctx *gin.Context, input *DTO.UpdateProjectValid) (common.Empty, error)
	DelProjectsFlow(ctx *gin.Context, input *DTO.DelProjectValid) (common.Empty, error)
	//ChangeProjectOwner(ctx *gin.Context, input *DTO.DelProjectValid) (common.Empty, error)
}

type ProjectMembersUseCase interface {
	GetProjectMembersFlow(ctx *gin.Context, input *DTO.MembersProjectIdValid) (DTO.AnswerProjectID, error)
	PostProjectMembersFlow(ctx *gin.Context, input *DTO.MembersAddIdValid) (DTO.AnswerMembers, error)
	PutProjectMembersFlow(ctx *gin.Context, input *DTO.MembersAddIdValid) (DTO.AnswerMembers, error)
	DelProjectMembersFlow(ctx *gin.Context, input *DTO.MembersAddIdValid) (DTO.AnswerMembers, error)
}
