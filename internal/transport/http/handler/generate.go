package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"go.uber.org/zap"

	appcontainer "scaffold/internal/app"
	"scaffold/internal/app/generator"
	"scaffold/internal/pkg/logger"
	"scaffold/internal/pkg/resp"
	ziputil "scaffold/internal/pkg/zip"
)

type generateRequest struct {
	TemplateID string                 `json:"template_id"`
	Version    string                 `json:"version"`
	Variables  map[string]interface{} `json:"variables"`
}

func Generate(ctx context.Context, c *app.RequestContext) {
	container := appcontainer.Container
	if container == nil || container.TemplateService == nil {
		resp.InternalError(c, "service not initialized")
		return
	}

	var req generateRequest
	if err := c.BindAndValidate(&req); err != nil {
		resp.BadRequest(c, "invalid request: "+err.Error())
		return
	}

	if req.TemplateID == "" {
		resp.BadRequest(c, "template_id is required")
		return
	}

	// Switch version if requested
	if req.Version != "" {
		if err := container.TemplateService.RefreshTemplateVersion(req.TemplateID, req.Version); err != nil {
			logger.Error("Failed to switch template version",
				zap.String("template_id", req.TemplateID),
				zap.String("version", req.Version),
				zap.Error(err))
			resp.InternalError(c, "failed to switch template version: "+err.Error())
			return
		}
	}

	tmplSource, err := container.TemplateService.GetTemplate(req.TemplateID)
	if err != nil {
		resp.NotFound(c, err.Error())
		return
	}

	// Variables are already map[string]interface{}, use directly
	vars := req.Variables
	if vars == nil {
		vars = make(map[string]interface{})
	}

	gen := generator.NewGenerator(
		tmplSource.LocalPath,
		tmplSource.Config,
		generator.WithVariables(vars),
	)

	files, err := gen.Generate()
	if err != nil {
		logger.Error("Code generation failed",
			zap.String("template_id", req.TemplateID),
			zap.Error(err))
		resp.InternalError(c, "code generation failed: "+err.Error())
		return
	}

	// Create ZIP from generated files
	zipBuf, err := ziputil.CreateFromFiles(files)
	if err != nil {
		logger.Error("Failed to create ZIP",
			zap.String("template_id", req.TemplateID),
			zap.Error(err))
		resp.InternalError(c, "failed to create archive: "+err.Error())
		return
	}

	// Determine filename
	projectName := "project"
	if name, ok := req.Variables["project_name"]; ok {
		if s, ok := name.(string); ok && s != "" {
			projectName = s
		}
	}

	// Send binary ZIP response
	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename="+projectName+".zip")
	c.Data(consts.StatusOK, "application/zip", zipBuf.Bytes())
}
