package main

import (
	"example/pdfgenerator/controllers"
	"example/pdfgenerator/initializers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.MigrateDB()
	initializers.InitMinioClient()
}

func main() {

	r := gin.Default()

	// Set up CORS middleware
	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	// r.POST("/generate", controllers.CreateOutput)
	// r.POST("/generate", controllers.CreateOutput)
	// r.GET("/generated-file", controllers.GeneratedOutput)
	// r.DELETE("/documents/:id", controllers.DeleteDocument)
	// r.GET("/documents/:id", controllers.GetDocumentById)
	// r.POST("/generate/:id", controllers.CreateDocument)

	r.POST("/upload-template", controllers.UploadTemplate)
	r.POST("/generate", controllers.CreateDocument, controllers.AutodocsLogs)
	r.GET("/documents", controllers.GetDocuments)
	r.GET("/templates", controllers.Templates)
	r.GET("/document-history", controllers.GetDocumentHistory)
	r.GET("/logs", controllers.AutodocsLogs)

	r.GET("/templates/preview/:refNumber", controllers.PreviewTemplate)
	r.GET("/documents/preview/:refNumber", controllers.PreviewDocument)

	r.DELETE("/templates/:refNumber", controllers.DeleteTemplate)
	r.DELETE("/documents/:refNumber", controllers.DeleteDocument)
	r.Run()

}
