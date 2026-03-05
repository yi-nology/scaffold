package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	appcontainer "scaffold/internal/app"
	"scaffold/internal/pkg/resp"
)

func Health(ctx context.Context, c *app.RequestContext) {
	resp.Success(c, map[string]interface{}{
		"status":  "ok",
		"version": "1.0.0",
	})
}

func ListTemplates(ctx context.Context, c *app.RequestContext) {
	container := appcontainer.Container
	if container == nil || container.TemplateService == nil {
		resp.InternalError(c, "service not initialized")
		return
	}

	templates := container.TemplateService.ListTemplates()
	resp.Success(c, templates)
}

func GetTemplate(ctx context.Context, c *app.RequestContext) {
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

	tmpl, err := container.TemplateService.GetTemplate(id)
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	resp.Success(c, tmpl)
}

func GetTemplateTags(ctx context.Context, c *app.RequestContext) {
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

	tags, err := container.TemplateService.GetTemplateTags(id)
	if err != nil {
		resp.InternalError(c, err.Error())
		return
	}

	resp.Success(c, tags)
}
