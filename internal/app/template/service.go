package template

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"scaffold/internal/conf"
	"scaffold/internal/pkg/logger"
	"scaffold/internal/repo/db/dao"
	"scaffold/internal/repo/db/model"
	"scaffold/internal/repo/redis/cache"
)

const ConfigFileName = "scaffold.yaml"

type TemplateSource struct {
	ID         string               `json:"id"`
	Repository string               `json:"repository"`
	LocalPath  string               `json:"local_path"`
	Config     *conf.TemplateConfig `json:"config"`
}

type Service struct {
	gitClient  *GitClient
	templates  map[string]*TemplateSource
	repository dao.TemplateRepository
	cache      *cache.TemplateCache
	dbEnabled  bool
}

type ServiceOption func(*Service)

func WithRepository(repo dao.TemplateRepository) ServiceOption {
	return func(s *Service) {
		s.repository = repo
		s.dbEnabled = repo != nil
	}
}

func WithCache(c *cache.TemplateCache) ServiceOption {
	return func(s *Service) {
		s.cache = c
	}
}

func NewService(cacheDir string, opts ...ServiceOption) *Service {
	s := &Service{
		gitClient: NewGitClient(cacheDir),
		templates: make(map[string]*TemplateSource),
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Service) AddTemplate(id, repoURL string) error {
	localPath, err := s.gitClient.Clone(repoURL)
	if err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}

	templateConfig, err := s.loadConfig(localPath)
	if err != nil {
		return fmt.Errorf("failed to load template config: %w", err)
	}

	s.templates[id] = &TemplateSource{
		ID:         id,
		Repository: repoURL,
		LocalPath:  localPath,
		Config:     templateConfig,
	}

	if s.dbEnabled {
		if err := s.syncToDB(id, repoURL, localPath, templateConfig); err != nil {
			logger.Warn("Failed to sync template to database",
				zap.String("id", id), zap.Error(err))
		}
	}

	return nil
}

func (s *Service) AddLocalTemplate(id, localPath string) error {
	templateConfig, err := s.loadConfig(localPath)
	if err != nil {
		return fmt.Errorf("failed to load template config: %w", err)
	}

	s.templates[id] = &TemplateSource{
		ID:        id,
		LocalPath: localPath,
		Config:    templateConfig,
	}

	if s.dbEnabled {
		if err := s.syncToDB(id, "", localPath, templateConfig); err != nil {
			logger.Warn("Failed to sync template to database",
				zap.String("id", id), zap.Error(err))
		}
	}

	return nil
}

func (s *Service) AddTemplateWithVersion(id, repoURL, version string) error {
	localPath, err := s.gitClient.CloneWithRef(repoURL, version)
	if err != nil {
		return fmt.Errorf("failed to clone template: %w", err)
	}

	templateConfig, err := s.loadConfig(localPath)
	if err != nil {
		return fmt.Errorf("failed to load template config: %w", err)
	}

	s.templates[id] = &TemplateSource{
		ID:         id,
		Repository: repoURL,
		LocalPath:  localPath,
		Config:     templateConfig,
	}

	if s.dbEnabled {
		if err := s.syncToDB(id, repoURL, localPath, templateConfig); err != nil {
			logger.Warn("Failed to sync template to database",
				zap.String("id", id), zap.Error(err))
		}
	}

	return nil
}

func (s *Service) RefreshTemplateVersion(id, version string) error {
	source, exists := s.templates[id]
	if !exists {
		return fmt.Errorf("template not found: %s", id)
	}

	if source.Repository == "" {
		return fmt.Errorf("template %s is not a remote template", id)
	}

	localPath, err := s.gitClient.CloneWithRef(source.Repository, version)
	if err != nil {
		return fmt.Errorf("failed to clone template version: %w", err)
	}

	templateConfig, err := s.loadConfig(localPath)
	if err != nil {
		return fmt.Errorf("failed to load template config: %w", err)
	}

	source.LocalPath = localPath
	source.Config = templateConfig

	return nil
}

func (s *Service) GetTemplate(id string) (*TemplateSource, error) {
	source, exists := s.templates[id]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", id)
	}
	return source, nil
}

func (s *Service) GetTemplateTags(id string) ([]cache.TagInfo, error) {
	source, exists := s.templates[id]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", id)
	}

	if source.Repository == "" {
		return []cache.TagInfo{}, nil
	}

	// Use cache if available
	if s.cache != nil {
		repoURL := source.Repository
		return s.cache.GetTags(id, func() ([]cache.TagInfo, error) {
			tags, err := s.gitClient.ListTagsWithAnnotations(repoURL)
			if err != nil {
				return nil, err
			}
			result := make([]cache.TagInfo, len(tags))
			for i, t := range tags {
				result[i] = cache.TagInfo{Name: t.Name, Message: t.Message}
			}
			return result, nil
		})
	}

	// Direct fetch without cache
	tags, err := s.gitClient.ListTagsWithAnnotations(source.Repository)
	if err != nil {
		return nil, err
	}
	result := make([]cache.TagInfo, len(tags))
	for i, t := range tags {
		result[i] = cache.TagInfo{Name: t.Name, Message: t.Message}
	}
	return result, nil
}

func (s *Service) DeleteTemplate(id string) error {
	delete(s.templates, id)

	if s.dbEnabled && s.repository != nil {
		ctx := context.Background()
		if err := s.repository.Delete(ctx, id); err != nil {
			logger.Warn("Failed to delete template from database",
				zap.String("id", id), zap.Error(err))
		}
	}

	if s.cache != nil {
		s.cache.ClearTemplate(id)
	}

	return nil
}

func (s *Service) ListTemplates() []conf.TemplateMeta {
	var result []conf.TemplateMeta
	for id, source := range s.templates {
		meta := conf.TemplateMeta{
			ID:          id,
			Name:        source.Config.Name,
			Description: source.Config.Description,
			Author:      source.Config.Author,
			Tags:        source.Config.Tags,
			Repository:  source.Repository,
		}
		result = append(result, meta)
	}
	return result
}

func (s *Service) LoadFromDatabase() error {
	if !s.dbEnabled || s.repository == nil {
		return nil
	}

	ctx := context.Background()
	models, err := s.repository.List(ctx)
	if err != nil {
		return fmt.Errorf("failed to load templates from database: %w", err)
	}

	for _, m := range models {
		if m.LocalPath == "" {
			logger.Warn("Template has no local path, skipping", zap.String("id", m.ID))
			continue
		}

		if _, err := os.Stat(m.LocalPath); os.IsNotExist(err) {
			logger.Warn("Template local path not exists, skipping",
				zap.String("id", m.ID), zap.String("path", m.LocalPath))
			continue
		}

		templateConfig, err := s.loadConfig(m.LocalPath)
		if err != nil {
			logger.Warn("Failed to load config for template, skipping",
				zap.String("id", m.ID), zap.Error(err))
			continue
		}

		s.templates[m.ID] = &TemplateSource{
			ID:         m.ID,
			Repository: m.Repository,
			LocalPath:  m.LocalPath,
			Config:     templateConfig,
		}

		logger.Info("Loaded template from database", zap.String("id", m.ID))
	}

	return nil
}

func (s *Service) GetGitClient() *GitClient {
	return s.gitClient
}

func (s *Service) syncToDB(id, repoURL, localPath string, cfg *conf.TemplateConfig) error {
	if s.repository == nil {
		return nil
	}

	ctx := context.Background()
	m := &model.TemplateModel{
		ID:          id,
		Name:        cfg.Name,
		Description: cfg.Description,
		Author:      cfg.Author,
		Version:     cfg.Version,
		Repository:  repoURL,
		LocalPath:   localPath,
		Tags:        cfg.Tags,
	}

	exists, err := s.repository.Exists(ctx, id)
	if err != nil {
		return err
	}

	if exists {
		return s.repository.Update(ctx, m)
	}
	return s.repository.Create(ctx, m)
}

func (s *Service) loadConfig(templatePath string) (*conf.TemplateConfig, error) {
	configPath := filepath.Join(templatePath, ConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var templateConfig conf.TemplateConfig
	if err := yaml.Unmarshal(data, &templateConfig); err != nil {
		return nil, err
	}

	if len(templateConfig.Ignore) == 0 {
		templateConfig.Ignore = []string{".git", "scaffold.yaml"}
	}

	return &templateConfig, nil
}
