package controllers

import (
	"anubis/app/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerProjects struct {
	projectUC interfaces.ProjectsUseCase
}

func NewControllerProjects(projectUC interfaces.ProjectsUseCase) *ControllerProjects {
	return &ControllerProjects{projectUC: projectUC}
}

func (c *ControllerProjects) HandlerGETProjects(ctx *gin.Context) {
	handlers.UriHandler(ctx, c.projectUC.GetProjectsFlow)
}

func (c *ControllerProjects) HandlerPOSTProjects(ctx *gin.Context) {
	handlers.JsonHandler(ctx, c.projectUC.PostProjectsFlow)
}

func (c *ControllerProjects) HandlerPUTProjects(ctx *gin.Context) {
	handlers.UriJsonHandler(ctx, c.projectUC.UpdateProjectsFlow)
}

func (c *ControllerProjects) HandlerDELProjects(ctx *gin.Context) {
	handlers.UriHandler(ctx, c.projectUC.DelProjectsFlow)
}
