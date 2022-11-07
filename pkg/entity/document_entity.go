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
	TemplateID  string
	Template    Template
	Fields      []DocumentField
	Stage       int
	VerifierID  string `gorm:"type:varchar(36)"`
	Verifier    User   `gorm:"foreignKey:VerifierID"`
	VerifiedAt  time.Time
	SignerID    string `gorm:"type:varchar(36)"`
	Signer      User   `gorm:"foreignKey:SignerID"`
	SignedAt    time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type DocumentField struct {
	gorm.Model
	DocumentID      string `gorm:"type:varchar(36)"`
	TemplateFieldID uint
	TemplateField   TemplateField
	Value           string
}

type Template struct {
	gorm.Model
	Name   string
	Path   string
	Fields []TemplateField
}

type TemplateField struct {
	gorm.Model
	TemplateID uint
	Key        string
}
