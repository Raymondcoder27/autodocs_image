package services

import (
	"example/pdfgenerator/initializers"
	"example/pdfgenerator/models"

	"gorm.io/gorm"
)

var DB *gorm.DB // Assume this is initialized somewhere

func SavePDF(pdf models.Document) error {
	return initializers.DB.Create(&pdf).Error
}

func SaveTemplate(template *models.Template) error {
	return initializers.DB.Create(template).Error
}

// func UpdateDbDocumentRecord(body models.Document) error {
//create document variable
// document := models.Document{
// 	MinioPdfObjectName: body.MinioPdfObjectName,
// 	URL: body.URL,
// }

// return initializers.DB.
// }
