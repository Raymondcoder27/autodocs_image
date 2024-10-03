package services

import (
	"bytes"
	"errors"
	"example/pdfgenerator/initializers"
	"example/pdfgenerator/models"

	// "encoding/base64"
	"html/template"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func GeneratePDF(templateBytes []byte, data map[string]interface{}) ([]byte, error) {
	// Parse the HTML template
	tmpl, err := template.New("upload").Parse(string(templateBytes))
	if err != nil {
		return nil, err
	}

	// Create a buffer to store the filled template
	var filledTemplate bytes.Buffer

	// Execute the template with the JSON data, storing the result in the buffer
	if err := tmpl.Execute(&filledTemplate, data); err != nil {
		return nil, err
	}

	// Initialize a new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	// Add a new page to the PDF generator with the filled template content
	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(filledTemplate.Bytes())))
	if err := pdfg.Create(); err != nil {
		return nil, err
	}

	// Encode the PDF bytes to base64
	// pdfBase64 := base64.StdEncoding.EncodeToString(pdfg.Bytes())
	// pdfBase64 :=
	return pdfg.Bytes(), nil
}

func DeleteDocumentByRefNumber(refNumber string) error {
	var document models.Document

	// Find the document by refNumber
	if err := initializers.DB.Where("ref_number = ?", refNumber).First(&document).Error; err != nil {
		return errors.New("document not found")
	}

	//delete the document from minio
	err := DeleteFile("pdfs", document.ID)
	if err != nil {
		return errors.New("failed to delete document from storage: " + err.Error())
	}

	// Delete the document
	if err := initializers.DB.Delete(&document).Error; err != nil {
		return err
	}

	return nil
}

func DeleteTemplateByRefNumber(refNumber string) error {
	var template models.Template

	// Find the document by refNumber
	if err := initializers.DB.Where("ref_number = ?", refNumber).First(&template).Error; err != nil {
		return errors.New("template not found")
	}

	//delete the file from minio
	err := DeleteFile("templates", template.ID)
	if err != nil {
		return errors.New("failed to delete template from storage: " + err.Error())
	}
	// Delete the document
	if err := initializers.DB.Delete(&template).Error; err != nil {
		return err
	}

	return nil
}
