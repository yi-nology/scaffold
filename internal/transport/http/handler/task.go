package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	appcontainer "scaffold/internal/app"
	"scaffold/internal/pkg/resp"
)

func GetTask(ctx context.Context, c *app.RequestContext) {
	container := appcontainer.Container
	if container == nil || container.TaskService == nil {
		resp.InternalError(c, "service not initialized")
		return
	}

	id := c.Param("id")
	if id == "" {
		resp.BadRequest(c, "task id is required")
		return
	}

	task, found := container.TaskService.GetTask(id)
	if !found {
		resp.NotFound(c, "task not found")
		return
	}

	resp.Success(c, task)
}

func ListTasks(ctx context.Context, c *app.RequestContext) {
	container := appcontainer.Container
	if container == nil || container.TaskService == nil {
		resp.InternalError(c, "service not initialized")
		return
	}

	tasks := container.TaskService.GetAllTasks()
	resp.Success(c, tasks)
}
