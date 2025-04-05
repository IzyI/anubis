package controllers

import (
	"anubis/app/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerProjectMembers struct {
	projectMembersUC interfaces.ProjectMembersUseCase
}

func NewControllerProjectMembers(projectMembersUC interfaces.ProjectMembersUseCase) *ControllerProjectMembers {
	return &ControllerProjectMembers{projectMembersUC: projectMembersUC}
}

func (c *ControllerProjectMembers) HandlerGETProjectMembers(ctx *gin.Context) {
	handlers.UriHandler(ctx, c.projectMembersUC.GetProjectMembersFlow)
}

func (c *ControllerProjectMembers) HandlerPOSTProjectMembers(ctx *gin.Context) {
	handlers.UriJsonHandler(ctx, c.projectMembersUC.PostProjectMembersFlow)
}

func (c *ControllerProjectMembers) HandlerPUTProjectMembers(ctx *gin.Context) {
	handlers.UriJsonHandler(ctx, c.projectMembersUC.PutProjectMembersFlow)
}

func (c *ControllerProjectMembers) HandlerDELProjectMembers(ctx *gin.Context) {
	handlers.UriJsonHandler(ctx, c.projectMembersUC.DelProjectMembersFlow)
}
