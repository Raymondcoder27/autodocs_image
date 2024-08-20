package services

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

var staticCounter = 1

func PDFFileName(templateFileName string) string {
	baseName := strings.TrimSuffix(filepath.Base(templateFileName), filepath.Ext(templateFileName))
	return fmt.Sprintf("%s.pdf", baseName)
}

func GenerateReferenceNumber() string {
	currentTime := time.Now()
	dateString := currentTime.Format("250831")
	// timeString := currentTime.Format("150309")
	staticCounter++

	return fmt.Sprintf("D%s-%04d", dateString, staticCounter)
}
