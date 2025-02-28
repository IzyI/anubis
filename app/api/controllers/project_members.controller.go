package controllers

import (
	"anubis/app/api/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerProjectMembers struct {
	projectMembersUC interfaces.ProjectMembersUseCase
}

func NewControllerProjectMembers(projectMembersUC interfaces.ProjectMembersUseCase) *ControllerProjectMembers {
	return &ControllerProjectMembers{projectMembersUC: projectMembersUC}
}

func (c *ControllerProjectMembers) HandlerGetProjectMembersFlow(ctx *gin.Context) {
	handlers.GetHandler(ctx, c.projectMembersUC.GetProjectMembersFlow)
}

func (c *ControllerProjectMembers) HandlerPOSTProjectMembersFlow(ctx *gin.Context) {
	handlers.PostHandler(ctx, c.projectMembersUC.PostProjectMembersFlow)
}

func (c *ControllerProjectMembers) HandlerPUTProjectMembersFlow(ctx *gin.Context) {
	handlers.PutHandler(ctx, c.projectMembersUC.PutProjectMembersFlow)
}

func (c *ControllerProjectMembers) HandlerDELProjectMembersFlow(ctx *gin.Context) {
	handlers.DeleteHandler(ctx, c.projectMembersUC.DelProjectMembersFlow)
}
