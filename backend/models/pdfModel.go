package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	ID           string `json:"id" gorm:"primaryKey"`
	DocumentName string `json:"documentName"`
	Description  string `json:"description"`
	TemplateId   string `json:"templateId"`
	// Data         string         `json:"pdf"`
	Status      string         `json:"requestStatus"`
	Method      string         `json:"requestMethod"`
	JsonPayload string         `json:"jsonPayload"`
	RefNumber   string         `json:"refNumber"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	// MinioPdfObjectName      string `json:"minioPdfObjectName"`
	// OriginalFileName        string `json:"originalTemplateName"`
	// URL                     string `json:"pdfFileUrl"`
}

type Template struct {
	ID        string         `json:"id"`
	Name      string         `json:"templateName"`
	RefNumber string         `json:"refNumber"`
	FileName  string         `json:"fileName"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Status    string         `json:"requestStatus"`
	Method    string         `json:"requestMethod"`
}

type Logs struct {
	ID                  string `json:"id"`
	DocumentName        string `json:"documentName"`
	DocumentDescription string `json:"documentDescription"`
	LogDescription      string `json:"logDescription"`
	TemplateId          string `json:"templateId"`
	// Data         string         `json:"pdf"`
	Status      string         `json:"requestStatus"`
	Method      string         `json:"requestMethod"`
	JsonPayload string         `json:"jsonPayload"`
	RefNumber   string         `json:"refNumber"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}
