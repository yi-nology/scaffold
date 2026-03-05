package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"scaffold/internal/repo/db/model"
)

// TemplateRepository interface for template data access
type TemplateRepository interface {
	Create(ctx context.Context, template *model.TemplateModel) error
	Update(ctx context.Context, template *model.TemplateModel) error
	GetByID(ctx context.Context, id string) (*model.TemplateModel, error)
	List(ctx context.Context) ([]*model.TemplateModel, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

type gormTemplateRepository struct {
	db *gorm.DB
}

func NewGormTemplateRepository(db *gorm.DB) TemplateRepository {
	return &gormTemplateRepository{db: db}
}

func (r *gormTemplateRepository) Create(ctx context.Context, template *model.TemplateModel) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *gormTemplateRepository) Update(ctx context.Context, template *model.TemplateModel) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *gormTemplateRepository) GetByID(ctx context.Context, id string) (*model.TemplateModel, error) {
	var tmpl model.TemplateModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&tmpl).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tmpl, nil
}

func (r *gormTemplateRepository) List(ctx context.Context) ([]*model.TemplateModel, error) {
	var templates []*model.TemplateModel
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func (r *gormTemplateRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.TemplateModel{}).Error
}

func (r *gormTemplateRepository) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.TemplateModel{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
