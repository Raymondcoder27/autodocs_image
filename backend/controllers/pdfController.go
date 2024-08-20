package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"example/pdfgenerator/initializers"
	"example/pdfgenerator/models"
	"example/pdfgenerator/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Code      string    `json:"code"`
	Timestamp time.Time `json:"currentTimestamp"`
}

type PDFResponse struct {
	ID           string    `json:"id"`
	DocumentName string    `json:"documentName"`
	TemplateId   string    `json:"templateId"`
	RefNumber    string    `json:"refNumber"`
	CreatedAt    time.Time `json:"createdAt"`
}

type PDFGenerationResponse struct {
	RefNumber string    `json:"refNumber"`
	CreatedAt time.Time `json:"createdAt"`
}

type TemplateUploadResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"templateName"`
	RefNumber string    `json:"refNumber"`
	FileName  string    `json:"fileName"`
	CreatedAt time.Time `json:"createdAt"`
}

type GenerateRequest struct {
	RefNumber   string                 `json:"refNumber"`
	Description string                 `json:"description"`
	Data        map[string]interface{} `json:"data"`
}

// UploadTemplate handles uploading an HTML template to MinIO
func UploadTemplate(c *gin.Context) {
	refNumber := services.GenerateReferenceNumber()

	file, _, err := c.Request.FormFile("template")
	templateName := c.PostForm("name")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to retrieve file: " + err.Error()})
		return
	}
	defer file.Close()

	templateBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading file: " + err.Error()})
		return
	}

	id := uuid.New().String()
	objectName := id

	template := models.Template{
		ID:        id,
		Name:      templateName,
		RefNumber: refNumber,
		FileName:  objectName,
		CreatedAt: time.Now(),
	}

	templateReader := bytes.NewReader(templateBytes)
	if err := services.UploadTemplate("templates", objectName, templateReader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading template file: " + err.Error()})
		return
	}

	if err := services.SaveTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving template metadata: " + err.Error()})
		return
	}

	// c.IndentedJSON(http.StatusOK, template)
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": template, "time": template.CreatedAt})
}

// CreateDocument generates a PDF using a stored template and JSON data
func CreateDocument(c *gin.Context) {
	id := uuid.New().String()
	// refNumber := c.PostForm("refNumber")
	// jsonData := c.PostForm("data")

	var request GenerateRequest

	// Bind the JSON request to the struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	var template models.Template
	if err := initializers.DB.First(&template, "ref_number = ?", request.RefNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Template not found for refNumber: " + request.RefNumber})
		return
	}

	templateId := template.FileName
	templateKey := templateId
	templateBytes, err := services.DownloadFile("templates", templateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching template: " + err.Error()})
		return
	}

	// Convert the map to a JSON string
	jsonString, err := json.Marshal(request.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to convert data to JSON string: " + err.Error()})
		return
	}

	data, err := services.DecodeJSON(string(jsonString))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON data: " + err.Error()})
		return
	}

	pdfBytes, err := services.GeneratePDF(templateBytes, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating PDF: " + err.Error()})
		return
	}

	// originalFileName := template.Name
	// documentName := services.PDFFileName(originalFileName)
	objectName := id
	fileReader := bytes.NewReader(pdfBytes)

	if err := services.UploadFile("pdfs", objectName, fileReader); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading PDF: " + err.Error()})
		return
	}

	storageKey := services.GenerateReferenceNumber()
	document := models.Document{
		ID:           id,
		DocumentName: id,
		Description:  request.Description,
		TemplateId:   templateId,
		RefNumber:    storageKey,
		CreatedAt:    time.Now(),
	}

	if err := initializers.DB.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		return
	}

	pdfGenerationResponse := PDFGenerationResponse{
		// ID:           document.ID,
		// DocumentName: document.DocumentName,
		// TemplateId:   document.TemplateId,
		RefNumber: document.RefNumber,
		CreatedAt: document.CreatedAt,
	}

	// c.IndentedJSON(http.StatusOK, pdfGenerationResponse)
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": pdfGenerationResponse, "timestamp": pdfGenerationResponse.CreatedAt})
}

// GetDocuments retrieves all documents
func GetDocuments(c *gin.Context) {
	var documents []models.Document
	if err := initializers.DB.Find(&documents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching documents"})
		return
	}
	currentTime := time.Now()
	// c.IndentedJSON(http.StatusOK, documents)
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": documents, "timestamp": currentTime})
}

// PreviewDocument returns the PDF for a given document refNumber
func PreviewDocument(c *gin.Context) {
	refNo := c.Param("refNumber")

	var document models.Document
	if err := initializers.DB.Where("ref_number = ?", refNo).First(&document).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Document not found"})
		return
	}

	objectName := document.ID

	pdfBytes, err := services.DownloadFile("pdfs", objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching PDF: " + err.Error()})
		return
	}

	// Encode the PDF bytes to base64
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)

	// c.JSON(http.StatusOK, pdfBase64)
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": pdfBase64, "timestamp": document.CreatedAt})

}

// PreviewTemplate returns the template file content
func PreviewTemplate(c *gin.Context) {
	refNo := c.Param("refNumber")

	var template models.Template
	if err := initializers.DB.Where("ref_number = ?", refNo).First(&template).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Template not found"})
		return
	}

	objectName := template.ID
	templateBytes, err := services.DownloadFile("templates", objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching template: " + err.Error()})
		return
	}

	// c.Data(http.StatusOK, "text/html", templateBytes)
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": templateBytes, "timestamp": template.CreatedAt})

}

// DeleteDocument deletes a document by refNumber
func DeleteDocument(c *gin.Context) {
	refNumber := c.Param("refNumber")

	err := services.DeleteDocumentByRefNumber(refNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}
	currentTime := time.Now()
	// c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "message": "Document deleted successfully", "timestamp": currentTime})
}

// DeleteTemplate deletes a template by refNumber
func DeleteTemplate(c *gin.Context) {
	refNumber := c.Param("refNumber")

	err := services.DeleteTemplateByRefNumber(refNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template not found"})
		return
	}

	currentTime := time.Now()
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Template deleted successfully", "timestamp": currentTime})
}

// Templates retrieves all templates
func Templates(c *gin.Context) {
	var templates []models.Template
	if err := initializers.DB.Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching templates"})
		return
	}
	currentTime := time.Now()
	// c.IndentedJSON(http.StatusOK, templates)
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": templates, "timestamp": currentTime})
}
