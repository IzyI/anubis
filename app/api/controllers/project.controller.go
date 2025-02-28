package controllers

import (
	"anubis/app/api/api/interfaces"
	"anubis/app/core/handlers"
	"github.com/gin-gonic/gin"
)

type ControllerProjects struct {
	projectUC interfaces.ProjectsUseCase
}

func NewControllerProjects(projectUC interfaces.ProjectsUseCase) *ControllerProjects {
	return &ControllerProjects{projectUC: projectUC}
}

func (c *ControllerProjects) HandlerGetProjectsFlow(ctx *gin.Context) {
	handlers.GetHandler(ctx, c.projectUC.GetProjectsFlow)
}

func (c *ControllerProjects) HandlerPostProjectsFlow(ctx *gin.Context) {
	handlers.PostHandler(ctx, c.projectUC.PostProjectsFlow)
}

func (c *ControllerProjects) HandlerPutProjectsFlow(ctx *gin.Context) {
	handlers.PutHandler(ctx, c.projectUC.PutProjectsFlow)
}

func (c *ControllerProjects) HandlerDelProjectsFlow(ctx *gin.Context) {
	handlers.DeleteHandler(ctx, c.projectUC.DelProjectsFlow)
}
