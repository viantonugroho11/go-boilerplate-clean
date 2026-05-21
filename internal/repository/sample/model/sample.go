package model

import (
	"time"

	entitysample "go-boilerplate-clean/internal/entity/sample"

	"gorm.io/gorm"
)

type Sample struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Code      string `gorm:"size:64;not null"`
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;not null"`
	Status    string `gorm:"size:32;not null;default:open"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Sample) TableName() string {
	return "samples"
}

func ToEntity(m *Sample) entitysample.Sample {
	if m == nil {
		return entitysample.Sample{}
	}
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		deletedAt = &t
	}
	return entitysample.Sample{
		ID:        m.ID,
		Code:      m.Code,
		Name:      m.Name,
		Email:     m.Email,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

func ToModel(s entitysample.Sample) Sample {
	return Sample{
		ID:        s.ID,
		Code:      s.Code,
		Name:      s.Name,
		Email:     s.Email,
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
