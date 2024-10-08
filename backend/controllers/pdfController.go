package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"example/pdfgenerator/initializers"
	"example/pdfgenerator/models"
	"example/pdfgenerator/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentPreviewResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"templateName"`
	RefNumber string    `json:"refNumber"`
	FileName  string    `json:"fileName"`
	CreatedAt time.Time `json:"created_at"`
	Method    string    `json:"requestMethod"`
	Status    string    `json:"requestStatus"`
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

type DeleteResponse struct {
	Status    string    `json:"responseStatus"`
	Method    string    `json:"responseMethod"`
	Code      string    `json:"code"`
	Timestamp time.Time `json:"currentTimestamp"`
}

type DocumentHistory struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
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

	//inserting post request into logs table
	if err := initializers.DB.Create(&models.Logs{
		ID:                  id,
		DocumentName:        id,
		JsonPayload:         "",
		Status:              "SUCCESS",
		Method:              "POST",
		DocumentDescription: "--",
		TemplateId:          id,
		RefNumber:           refNumber,
		CreatedAt:           time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
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

	currentTime := time.Now()

	var request GenerateRequest

	// Bind the JSON request to the struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		//inserting get request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "POST",
			LogDescription: "Invalid Request",
			TemplateId:     "",
			RefNumber:      "",
			CreatedAt:      currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           id,
			DocumentName: id,
			Description:  request.Description,
			TemplateId:   "",
			Status:       "FAILED",
			Method:       "POST",
			JsonPayload:  "",
			RefNumber:    request.RefNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}
		return
	}

	var template models.Template
	if err := initializers.DB.First(&template, "ref_number = ?", request.RefNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Template not found for refNumber: " + request.RefNumber})
		//inserting post request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:                  id,
			DocumentName:        id,
			JsonPayload:         "",
			Status:              "FAILED",
			Method:              "POST",
			DocumentDescription: "--",
			TemplateId:          "",
			RefNumber:           request.RefNumber,
			CreatedAt:           time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           id,
			DocumentName: id,
			Description:  request.Description,
			TemplateId:   "",
			Status:       "FAILED",
			Method:       "POST",
			JsonPayload:  "",
			RefNumber:    request.RefNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}
		return
	}

	templateId := template.FileName
	templateKey := templateId
	templateBytes, err := services.DownloadFile("templates", templateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching template: " + err.Error()})
		//inserting post request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:                  id,
			DocumentName:        id,
			JsonPayload:         "",
			Status:              "FAILED",
			Method:              "POST",
			DocumentDescription: "Error fetching template: " + err.Error(),
			TemplateId:          "",
			RefNumber:           request.RefNumber,
			CreatedAt:           time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           id,
			DocumentName: id,
			Description:  request.Description,
			TemplateId:   templateId,
			Status:       "FAILED",
			Method:       "POST",
			JsonPayload:  "",
			RefNumber:    request.RefNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}
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
		//inserting post request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:                  id,
			DocumentName:        id,
			JsonPayload:         string(jsonString),
			Status:              "FAILED",
			Method:              "POST",
			DocumentDescription: request.Description,
			TemplateId:          templateId,
			RefNumber:           request.RefNumber,
			CreatedAt:           time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           id,
			DocumentName: id,
			Description:  request.Description,
			TemplateId:   templateId,
			Status:       "FAILED",
			Method:       "POST",
			JsonPayload:  string(jsonString),
			RefNumber:    request.RefNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}
		return
	}

	objectName := id
	fileReader := bytes.NewReader(pdfBytes)

	if err := services.UploadFile("pdfs", objectName, fileReader); err != nil {

		//inserting post request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:                  id,
			DocumentName:        id,
			JsonPayload:         string(jsonString),
			Status:              "FAILED",
			Method:              "POST",
			DocumentDescription: request.Description,
			TemplateId:          templateId,
			RefNumber:           request.RefNumber,
			CreatedAt:           time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           id,
			DocumentName: id,
			Description:  request.Description,
			TemplateId:   templateId,
			Status:       "FAILED",
			Method:       "POST",
			JsonPayload:  string(jsonString),
			RefNumber:    request.RefNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading PDF: " + err.Error()})
		return
	}

	storageKey := services.GenerateReferenceNumber()
	document := models.Document{
		ID:           id,
		DocumentName: id,
		JsonPayload:  string(jsonString),
		Description:  request.Description,
		TemplateId:   templateId,
		RefNumber:    storageKey,
		CreatedAt:    time.Now(),
	}

	if err := initializers.DB.Create(&document).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		//inserting post request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:                  id,
			DocumentName:        id,
			JsonPayload:         string(jsonString),
			Status:              "FAILED",
			Method:              "POST",
			DocumentDescription: request.Description,
			TemplateId:          templateId,
			RefNumber:           request.RefNumber,
			CreatedAt:           time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           id,
			DocumentName: id,
			Description:  request.Description,
			TemplateId:   templateId,
			Status:       "FAILED",
			Method:       "POST",
			JsonPayload:  string(jsonString),
			RefNumber:    request.RefNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}
		return
	}

	//inserting post request into logs table
	if err := initializers.DB.Create(&models.Logs{
		ID:                  id,
		DocumentName:        id,
		JsonPayload:         string(jsonString),
		Status:              "SUCCESS",
		Method:              "POST",
		DocumentDescription: request.Description,
		TemplateId:          templateId,
		RefNumber:           storageKey,
		CreatedAt:           time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		return
	}

	pdfGenerationResponse := PDFGenerationResponse{
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

	//inserting get request into logs table
	// if err := initializers.DB.Create(&models.Logs{
	// 	ID:             uuid.New().String(),
	// 	DocumentName:   "",
	// 	JsonPayload:    "",
	// 	Status:         "SUCCESS",
	// 	Method:         "GET",
	// 	LogDescription: "Get all documents",
	// 	TemplateId:     "",
	// 	RefNumber:      "",
	// 	CreatedAt:      currentTime,
	// }).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
	// 	return
	// }

	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": documents, "timestamp": currentTime})
}

// PreviewDocument returns the PDF for a given document refNumber
func PreviewDocument(c *gin.Context) {
	refNo := c.Param("refNumber")

	var document models.Document
	if err := initializers.DB.Where("ref_number = ?", refNo).First(&document).Error; err != nil {

		//inserting get request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "GET",
			LogDescription: "Document not found",
			TemplateId:     "",
			RefNumber:      refNo,
			CreatedAt:      time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Document not found"})
		return
	}

	objectName := document.ID

	pdfBytes, err := services.DownloadFile("pdfs", objectName)
	if err != nil {
		//inserting get request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "GET",
			LogDescription: "Error fetching PDF: " + err.Error(),
			TemplateId:     "",
			RefNumber:      refNo,
			CreatedAt:      time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching PDF: " + err.Error()})
		return
	}

	// Encode the PDF bytes to base64
	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)

	//insert into logs table
	if err := initializers.DB.Create(&models.Logs{
		ID:                  document.ID,
		DocumentName:        document.ID,
		JsonPayload:         "",
		Status:              "SUCCESS",
		Method:              "GET",
		DocumentDescription: document.Description,
		TemplateId:          document.TemplateId,
		RefNumber:           refNo,
		CreatedAt:           time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		return
	}

	// c.JSON(http.StatusOK, pdfBase64)
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": pdfBase64, "timestamp": document.CreatedAt})
}

// PreviewTemplate returns the template file content
func PreviewTemplate(c *gin.Context) {
	refNo := c.Param("refNumber")

	var template models.Template
	if err := initializers.DB.Where("ref_number = ?", refNo).First(&template).Error; err != nil {
		//inserting get request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "GET",
			LogDescription: "Template not found",
			TemplateId:     "",
			RefNumber:      refNo,
			CreatedAt:      time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Template not found"})
		return
	}

	objectName := template.ID
	templateBytes, err := services.DownloadFile("templates", objectName)
	if err != nil {
		//inserting get request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "GET",
			LogDescription: "Error fetching template",
			TemplateId:     "",
			RefNumber:      refNo,
			CreatedAt:      time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching template: " + err.Error()})
		return
	}

	// c.Data(http.StatusOK, "text/html", templateBytes)
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": templateBytes, "timestamp": template.CreatedAt})

}

// DeleteDocument deletes a document by refNumber
func DeleteDocument(c *gin.Context) {
	refNumber := c.Param("refNumber")

	currentTime := time.Now()

	//find this document in the database
	var document models.Document
	if err := initializers.DB.Where("ref_number = ?", refNumber).First(&document).Error; err != nil {
		if err := initializers.DB.Create(&models.Logs{
			// ID: uuid.New().String(),
			ID:                  document.ID,
			DocumentName:        document.ID,
			DocumentDescription: document.Description,
			TemplateId:          document.TemplateId,
			JsonPayload:         "",
			Status:              "FAILED",
			Method:              "DELETE",
			LogDescription:      "Document not found",
			RefNumber:           refNumber,
			CreatedAt:           currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           document.ID,
			DocumentName: document.ID,
			Description:  document.Description,
			TemplateId:   document.TemplateId,
			Status:       "FAILED",
			Method:       "DELETE",
			JsonPayload:  "",
			RefNumber:    refNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}
		c.JSON(http.StatusNotFound, gin.H{"message": "Document not found"})
		return
	}

	//inserting delete request into logs table
	if err := initializers.DB.Create(&models.Logs{
		// ID: uuid.New().String(),
		ID:                  document.ID,
		DocumentName:        document.ID,
		DocumentDescription: document.Description,
		TemplateId:          document.TemplateId,
		JsonPayload:         "",
		Status:              "SUCCESS",
		Method:              "DELETE",
		LogDescription:      "Document deleted successfully",
		RefNumber:           refNumber,
		CreatedAt:           currentTime,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		return
	}

	err := services.DeleteDocumentByRefNumber(refNumber)
	if err != nil {
		if err := initializers.DB.Create(&models.Logs{
			// ID: uuid.New().String(),
			ID:                  document.ID,
			DocumentName:        document.ID,
			DocumentDescription: document.Description,
			TemplateId:          document.TemplateId,
			JsonPayload:         "",
			Status:              "FAILED",
			Method:              "DELETE",
			LogDescription:      "Failed to delete document",
			RefNumber:           refNumber,
			CreatedAt:           currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	response := DeleteResponse{
		Status:    "SUCCESS",
		Method:    "DELETE",
		Code:      "200",
		Timestamp: currentTime,
	}
	c.IndentedJSON(http.StatusOK, response)
}

// DeleteTemplate deletes a template by refNumber
func DeleteTemplate(c *gin.Context) {
	refNumber := c.Param("refNumber")

	//find this template in the database
	var template models.Template
	if err := initializers.DB.Where("ref_number = ?", refNumber).First(&template).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Template not found"})
		return
	}

	//inserting delete request into logs table
	if err := initializers.DB.Create(&models.Logs{
		// ID: uuid.New().String(),
		ID:             uuid.New().String(),
		DocumentName:   "",
		JsonPayload:    "",
		Status:         "SUCCESS",
		Method:         "DELETE",
		LogDescription: "Template deleted successfully",
		TemplateId:     template.ID,
		RefNumber:      refNumber,
		CreatedAt:      time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		return
	}

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

func GetDocumentHistory(c *gin.Context) {
	var history []DocumentHistory

	// Calculate the start of the current week (most recent Sunday)
	now := time.Now()
	weekday := int(now.Weekday())
	startOfWeek := now.AddDate(0, 0, -weekday)

	// Group by creation date and count documents within the current week
	err := initializers.DB.Table("documents").
		Select("TO_CHAR(created_at AT TIME ZONE 'UTC', 'FMDay') as date, COUNT(*) as count").
		Where("created_at >= ?", startOfWeek).
		Group("TO_CHAR(created_at AT TIME ZONE 'UTC', 'FMDay')").
		Scan(&history).Error

	if err != nil {
		log.Println("Error fetching document history:", err)

		//inserting get request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "GET",
			LogDescription: "Error fetching document history",
			TemplateId:     "",
			RefNumber:      "",
			CreatedAt:      now,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching document history"})
		return
	}

	// Initialize counts for each day of the week
	dayCounts := map[string]int{
		"Sunday":    0,
		"Monday":    0,
		"Tuesday":   0,
		"Wednesday": 0,
		"Thursday":  0,
		"Friday":    0,
		"Saturday":  0,
	}

	// Populate counts from the database results
	for _, record := range history {
		dayCounts[record.Date] = record.Count
	}

	// Get the current day of the week
	currentDay := now.Weekday().String()

	// Create the final response in the desired order, with the current day last
	var response []map[string]interface{}
	orderedDays := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
	for _, day := range orderedDays {
		if day != currentDay {
			response = append(response, map[string]interface{}{
				"date":  day,
				"count": dayCounts[day],
			})
		}
	}
	// Add the current day last
	response = append(response, map[string]interface{}{
		"date":  currentDay,
		"count": dayCounts[currentDay],
	})

	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": response})
}

func AutodocsLogs(c *gin.Context) {
	var logs []models.Logs
	if err := initializers.DB.Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching logs"})
		return
	}
	currentTime := time.Now()
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": logs, "timestamp": currentTime})
}

func DeleteAllLogs(c *gin.Context) {
	err := initializers.DB.Exec("DELETE FROM logs").Error
	if err != nil {

		//inserting delete request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "DELETE",
			LogDescription: "Error deleting logs",
			TemplateId:     "",
			RefNumber:      "",
			CreatedAt:      time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting logs"})
		return
	}
	currentTime := time.Now()
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "Logs deleted successfully", "timestamp": currentTime})
}

// GetMetrics retrieves metrics based on the provided date range
func GetRangeMetrics(c *gin.Context) {
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")

	var totalTemplates int64
	var totalDocuments int64
	var failedGenerations int64

	// Parse the start and end dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {

		//inserting get request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "GET",
			LogDescription: "Invalid start date",
			TemplateId:     "",
			RefNumber:      "",
			CreatedAt:      time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid start date"})
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {

		//inserting get request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "GET",
			LogDescription: "Invalid end date",
			TemplateId:     "",
			RefNumber:      "",
			CreatedAt:      time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid end date"})
		return
	}

	// Count total templates within the date range
	if err := initializers.DB.Model(&models.Template{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&totalTemplates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching templates count"})
		return
	}

	// Count total documents within the date range
	if err := initializers.DB.Model(&models.Document{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&totalDocuments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching documents count"})
		return
	}

	// Count failed generations within the date range
	// if err := initializers.DB.Model(&models.Logs{}).
	// 	Where("created_at BETWEEN ? AND ? AND status = ?", start, end, "FAILED").
	// 	Count(&failedGenerations).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching failed generations count"})
	// 	return
	// }

	//Count failed generations from failed generations table
	if err := initializers.DB.Model(&models.FailedGenerations{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&failedGenerations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching failed generations count"})
		return
	}

	//calculate generation rate and failure rate
	var generationRate float64
	var failureRate float64
	totalDays := end.Sub(start).Hours() / 24

	if totalDays > 0 {
		generationRate = float64(totalDocuments) / totalDays
	}

	if totalDocuments > 0 {
		failureRate = (float64(failedGenerations) / float64(totalDocuments)) * 100
	}

	// Prepare the response
	response := gin.H{
		"totalTemplates":    totalTemplates,
		"totalDocuments":    totalDocuments,
		"failedGenerations": failedGenerations,
		"generationRate":    generationRate,
		"failureRate":       failureRate,
		"timestamp":         time.Now(),
	}

	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": response})
}

// failed Generations
func GetFailedGenerations(c *gin.Context) {
	var failedGenerations []models.FailedGenerations
	if err := initializers.DB.Find(&failedGenerations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching failed generations"})
		return
	}
	currentTime := time.Now()
	c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": failedGenerations, "timestamp": currentTime})
}

func HtmlBeforePDF(c *gin.Context) {
	id := uuid.New().String()
	// refNumber := c.PostForm("refNumber")
	// jsonData := c.PostForm("data")

	currentTime := time.Now()

	var request GenerateRequest

	// Bind the JSON request to the struct
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		//inserting get request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:             uuid.New().String(),
			DocumentName:   "",
			JsonPayload:    "",
			Status:         "FAILED",
			Method:         "POST",
			LogDescription: "Invalid Request",
			TemplateId:     "",
			RefNumber:      "",
			CreatedAt:      currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           id,
			DocumentName: id,
			Description:  request.Description,
			TemplateId:   "",
			Status:       "FAILED",
			Method:       "POST",
			JsonPayload:  "",
			RefNumber:    request.RefNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}
		return
	}

	var template models.Template
	if err := initializers.DB.First(&template, "ref_number = ?", request.RefNumber).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Template not found for refNumber: " + request.RefNumber})
		//inserting post request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:                  id,
			DocumentName:        id,
			JsonPayload:         "",
			Status:              "FAILED",
			Method:              "POST",
			DocumentDescription: "--",
			TemplateId:          "",
			RefNumber:           request.RefNumber,
			CreatedAt:           time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           id,
			DocumentName: id,
			Description:  request.Description,
			TemplateId:   "",
			Status:       "FAILED",
			Method:       "POST",
			JsonPayload:  "",
			RefNumber:    request.RefNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}
		return
	}

	templateId := template.FileName
	templateKey := templateId
	templateBytes, err := services.DownloadFile("templates", templateKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching template: " + err.Error()})
		//inserting post request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:                  id,
			DocumentName:        id,
			JsonPayload:         "",
			Status:              "FAILED",
			Method:              "POST",
			DocumentDescription: "Error fetching template: " + err.Error(),
			TemplateId:          "",
			RefNumber:           request.RefNumber,
			CreatedAt:           time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}

		//insert into failed generations table
		if err := initializers.DB.Create(&models.FailedGenerations{
			ID:           id,
			DocumentName: id,
			Description:  request.Description,
			TemplateId:   templateId,
			Status:       "FAILED",
			Method:       "POST",
			JsonPayload:  "",
			RefNumber:    request.RefNumber,
			CreatedAt:    currentTime,
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
		}
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

	htmlBeforePDF, err := services.GeneratePDF2(templateBytes, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating PDF: " + err.Error()})
		//inserting post request into logs table
		if err := initializers.DB.Create(&models.Logs{
			ID:                  id,
			DocumentName:        id,
			JsonPayload:         string(jsonString),
			Status:              "FAILED",
			Method:              "POST",
			DocumentDescription: request.Description,
			TemplateId:          templateId,
			RefNumber:           request.RefNumber,
			CreatedAt:           time.Now(),
		}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving document metadata in database: " + err.Error()})
			return
		}
	}

	fmt.Printf("---------------------------------------------")

	// fmt.Print(htmlBeforePDF)
	// log.Print(htmlBeforePDF)

	fileData := []byte(htmlBeforePDF)
	err = os.WriteFile("file.txt", fileData, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to write file: " + err.Error()})
		return
	}

	fmt.Printf("File generated and written successfully\n")
	// c.IndentedJSON(http.StatusOK, gin.H{"code": 200, "data": htmlBeforePDF, "timestamp": currentTime})
}
