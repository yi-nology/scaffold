package app

import (
	"scaffold/internal/app/generator"
	"scaffold/internal/app/task"
	tmplservice "scaffold/internal/app/template"
)

// Container holds all service instances, initialized by bootstrap
var Container *ServiceContainer

type ServiceContainer struct {
	TemplateService  *tmplservice.Service
	GeneratorService *generator.Generator // Generator is created per-request, this is nil
	TaskService      *task.Service
}
