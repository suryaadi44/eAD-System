package entity

import (
	"gorm.io/gorm"
	"time"
)

type Document struct {
	ID          string `gorm:"primaryKey; type:varchar(36)"`
	Register    string `gorm:"type:varchar(255);not null;uniqueIndex"`
	ApplicantID string `gorm:"type:varchar(36);not null"`
	Applicant   User   `gorm:"foreignKey:ApplicantID"`
	TemplateID  uint
	Template    Template
	Fields      DocumentFields
	StageID     int            `gorm:"type:int;default:1"`
	Stage       Stage          `gorm:"foreignKey:StageID"`
	VerifierID  string         `gorm:"type:varchar(36);default:null"`
	Verifier    User           `gorm:"foreignKey:VerifierID"`
	VerifiedAt  time.Time      `gorm:"type:datetime;default:null"`
	SignerID    string         `gorm:"type:varchar(36);default:null"`
	Signer      User           `gorm:"foreignKey:SignerID"`
	SignedAt    time.Time      `gorm:"type:datetime;default:null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type DocumentField struct {
	gorm.Model
	DocumentID      string `gorm:"type:varchar(36)"`
	TemplateFieldID uint
	TemplateField   TemplateField
	Value           string
}

type DocumentFields []DocumentField

type Stage struct {
	ID     int    `gorm:"primaryKey; type:int"`
	Status string `gorm:"type:varchar(255);not null;uniqueIndex"`
}

type Template struct {
	gorm.Model
	Name         string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Path         string `gorm:"type:varchar(255);not null;uniqueIndex"`
	MarginTop    uint
	MarginBottom uint
	MarginLeft   uint
	MarginRight  uint
	IsActive     bool `gorm:"default:true"`
	Fields       TemplateFields
}

type Templates []Template

type TemplateField struct {
	gorm.Model
	TemplateID uint
	Key        string
}

type TemplateFields []TemplateField
