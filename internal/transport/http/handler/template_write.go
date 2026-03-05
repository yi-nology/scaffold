package handler

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	appcontainer "scaffold/internal/app"
	"scaffold/internal/pkg/logger"
	"scaffold/internal/pkg/resp"
)

type addTemplateRequest struct {
	ID         string `json:"id"`
	Repository string `json:"repository"`
	Version    string `json:"version"`
}

func AddTemplate(ctx context.Context, c *app.RequestContext) {
	container := appcontainer.Container
	if container == nil || container.TemplateService == nil || container.TaskService == nil {
		resp.InternalError(c, "service not initialized")
		return
	}

	var req addTemplateRequest
	if err := c.BindAndValidate(&req); err != nil {
		resp.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if req.ID == "" || req.Repository == "" {
		resp.BadRequest(c, "id and repository are required")
		return
	}

	task := container.TaskService.CreateTask(req.ID)

	go func() {
		container.TaskService.UpdateTask(task.ID, "running", "正在克隆模板仓库...")

		var err error
		if req.Version != "" {
			err = container.TemplateService.AddTemplateWithVersion(req.ID, req.Repository, req.Version)
		} else {
			err = container.TemplateService.AddTemplate(req.ID, req.Repository)
		}

		if err != nil {
			logger.Error("Failed to add template",
				zap.String("id", req.ID),
				zap.Error(err))
			container.TaskService.UpdateTask(task.ID, "failed", "添加模板失败: "+err.Error())
			return
		}

		container.TaskService.UpdateTask(task.ID, "completed", "模板添加成功")
		logger.Info("Template added successfully", zap.String("id", req.ID))
	}()

	c.JSON(consts.StatusAccepted, map[string]interface{}{
		"code":    0,
		"message": "template addition started",
		"data": map[string]interface{}{
			"task_id": task.ID,
		},
	})
}

type deleteTemplateRequest struct {
	ID string `json:"id"`
}

func DeleteTemplate(ctx context.Context, c *app.RequestContext) {
	container := appcontainer.Container
	if container == nil || container.TemplateService == nil {
		resp.InternalError(c, "service not initialized")
		return
	}

	var req deleteTemplateRequest
	if err := c.BindAndValidate(&req); err != nil {
		resp.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if req.ID == "" {
		resp.BadRequest(c, "template id is required")
		return
	}

	if err := container.TemplateService.DeleteTemplate(req.ID); err != nil {
		resp.InternalError(c, "failed to delete template: "+err.Error())
		return
	}

	resp.Success(c, nil)
}

func RefreshTemplateVersion(ctx context.Context, c *app.RequestContext) {
	container := appcontainer.Container
	if container == nil || container.TemplateService == nil {
		resp.InternalError(c, "service not initialized")
		return
	}

	id := c.Param("id")
	if id == "" {
		resp.BadRequest(c, "template id is required")
		return
	}

	var body struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(c.Request.Body(), &body); err != nil {
		resp.BadRequest(c, "invalid request body")
		return
	}

	if body.Version == "" {
		resp.BadRequest(c, "version is required")
		return
	}

	if err := container.TemplateService.RefreshTemplateVersion(id, body.Version); err != nil {
		resp.InternalError(c, "failed to refresh template version: "+err.Error())
		return
	}

	resp.Success(c, nil)
}

func AddLocalTemplate(ctx context.Context, c *app.RequestContext) {
	container := appcontainer.Container
	if container == nil || container.TemplateService == nil {
		resp.InternalError(c, "service not initialized")
		return
	}

	var req struct {
		ID   string `json:"id"`
		Path string `json:"path"`
	}
	if err := c.BindAndValidate(&req); err != nil {
		resp.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if req.ID == "" || req.Path == "" {
		resp.BadRequest(c, "id and path are required")
		return
	}

	if err := container.TemplateService.AddLocalTemplate(req.ID, req.Path); err != nil {
		resp.InternalError(c, "failed to add local template: "+err.Error())
		return
	}

	resp.Success(c, nil)
}
