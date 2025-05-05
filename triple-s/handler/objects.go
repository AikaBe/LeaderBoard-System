package handler

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"triple-s/config"
)

func UploadObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Get the bucket and object key from URL
	// Strip the leading slash from the path and split by "/"
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	bucketName := parts[0]
	objectKey := parts[1]

	// Check if the bucket exists
	bucketPath := filepath.Join(StorageDir, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		http.Error(w, "Bucket does not exist", http.StatusNotFound)
		return
	}

	// Validate the object key
	if !isValidObjectKey(objectKey) {
		http.Error(w, "Invalid object key", http.StatusBadRequest)
		return
	}

	// Save the object content as a file
	objectFilePath := filepath.Join(bucketPath, objectKey)
	file, err := os.Create(objectFilePath)
	if err != nil {
		http.Error(w, "Error saving object", http.StatusInternalServerError)
		log.Printf("Error creating object file: %v", err)
		return
	}
	defer file.Close()

	// Copy the file content from the request body to the object file
	_, err = io.Copy(file, r.Body)
	if err != nil {
		http.Error(w, "Error uploading object", http.StatusInternalServerError)
		log.Printf("Error uploading object: %v", err)
		return
	}

	// Update object metadata in the CSV file
	err = updateObjectMetadata(bucketName, objectKey, r.Header.Get("Content-Type"), file)
	if err != nil {
		http.Error(w, "Error saving metadata", http.StatusInternalServerError)
		log.Printf("Error saving object metadata: %v", err)
		return
	}

	// host := "triple-s"
	// port := "9000"
	// imageURL := fmt.Sprintf("http://%s:%s/images/%s", host, port, objectKey)
	imageURL := fmt.Sprintf("http://localhost:9000/%s/%s", bucketName, objectKey)

	// Return JSON response with image URL
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(imageURL))
}

// Helper function to validate object key
func isValidObjectKey(objectKey string) bool {
	// Here you can implement your custom validation for object keys
	// For simplicity, we're allowing only alphanumeric characters, dashes, and underscores
	return objectKey != "" && len(objectKey) > 2 && len(objectKey) < 256
}

// Helper function to update the object metadata in the CSV file
func updateObjectMetadata(bucketName, objectKey, contentType string, file *os.File) error {
	// Open the CSV file for the bucket
	objectsCSVPath := filepath.Join(StorageDir, bucketName, "objects.csv")
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %w", err)
	}

	// Open the CSV file, or create it if it doesn't exist
	csvFile, err := os.OpenFile(objectsCSVPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("error opening objects metadata file: %w", err)
	}
	defer csvFile.Close()

	// Create a CSV writer
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write metadata record for the object
	record := []string{
		objectKey,
		contentType,
		fmt.Sprintf("%d", fileInfo.Size()),
		time.Now().Format(time.RFC3339),
	}

	err = writer.Write(record)
	if err != nil {
		return fmt.Errorf("error writing to objects CSV: %w", err)
	}

	return nil
}

func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract bucket and object name from URL path
	// URL structure: /bucketName/objectKey
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")
	if len(pathSegments) != 2 {
		http.Error(w, "Invalid URL structure. Should be /bucketName/objectKey", http.StatusBadRequest)
		return
	}

	bucketName := pathSegments[0]
	objectKey := pathSegments[1]

	// Verify bucket existence
	bucketDir := filepath.Join(config.GetStorage(), bucketName)
	if _, err := os.Stat(bucketDir); os.IsNotExist(err) {
		http.Error(w, "Bucket does not exist", http.StatusNotFound)
		return
	}

	// Verify object existence (check if the file exists in the bucket's "objects" folder)
	objectFilePath := filepath.Join(bucketDir, objectKey)
	if _, err := os.Stat(objectFilePath); os.IsNotExist(err) {
		http.Error(w, "Object does not exist", http.StatusNotFound)
		return
	}

	// Set appropriate Content-Type header (dynamically set based on file extension)
	contentType := getContentType(objectKey)
	w.Header().Set("Content-Type", contentType)

	// Set Content-Disposition to inline (i.e., display the content in the browser)
	// This ensures the browser tries to display the content directly instead of downloading it
	w.Header().Set("Content-Disposition", "inline; filename="+objectKey)

	// Open the object file and serve it
	file, err := os.Open(objectFilePath)
	if err != nil {
		http.Error(w, "Error opening object", http.StatusInternalServerError)
		log.Printf("Error opening object file: %v", err)
		return
	}
	defer file.Close()

	// Copy the file contents to the response
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error reading object content", http.StatusInternalServerError)
		log.Printf("Error reading object file: %v", err)
		return
	}

	// Respond with 200 OK and the content of the object
	w.WriteHeader(http.StatusOK)
}

// Helper function to get content type based on file extension
func getContentType(fileName string) string {
	// Get file extension
	ext := filepath.Ext(fileName)

	// Return the appropriate content type based on the extension
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".html":
		return "text/html"
	default:
		return "application/octet-stream" // Default for binary data
	}
}

// Respond with an error in XML format
func respondWithError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)
	w.Write([]byte(fmt.Sprintf("<error><message>%s</message></error>", message)))
}

// Respond with a success message in XML format
func respondWithSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<response><message>%s</message></response>", message)))
}

// Delete the object metadata from the CSV file
func removeObjectMetadata(metadataFilePath string, objectKey string) error {
	// Open the CSV file
	file, err := os.OpenFile(metadataFilePath, os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Find and remove the row for the objectKey
	var updatedRecords [][]string
	found := false
	for _, record := range records {
		if record[0] != objectKey { // Assuming the objectKey is the first column
			updatedRecords = append(updatedRecords, record)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("object metadata not found")
	}

	// Move the cursor back to the beginning of the file to overwrite it
	file.Seek(0, 0)

	// Write the updated records back to the CSV file
	writer := csv.NewWriter(file)
	err = writer.WriteAll(updatedRecords)
	if err != nil {
		return err
	}

	// Truncate the file to remove the old content that wasn't rewritten
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	return nil
}

// DeleteObjectHandler handles the deletion of an object
func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	// Extract bucket and object name from URL path
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/"), "/")

	// Ensure the URL contains exactly two parts (bucketName and objectKey)
	if len(pathSegments) != 2 {
		http.Error(w, "Invalid URL structure. Should be /bucketName/objectKey", http.StatusBadRequest)
		return
	}

	bucketName := pathSegments[0]
	objectKey := pathSegments[1]

	// Verify bucket existence
	bucketDir := filepath.Join(config.GetStorage(), bucketName)
	if _, err := os.Stat(bucketDir); os.IsNotExist(err) {
		// If bucket doesn't exist, return 404 with XML error message
		respondWithError(w, "Bucket not found", http.StatusNotFound)
		return
	}

	// Verify object existence
	objectFilePath := filepath.Join(bucketDir, objectKey)
	if _, err := os.Stat(objectFilePath); os.IsNotExist(err) {
		// If object doesn't exist, return 404 with XML error message
		respondWithError(w, "Object not found", http.StatusNotFound)
		return
	}

	// Delete the object file
	err := os.Remove(objectFilePath)
	if err != nil {
		// If there is an error deleting the object, return 500 with XML error message
		respondWithError(w, "Error deleting object", http.StatusInternalServerError)
		log.Printf("Error deleting object file: %v", err)
		return
	}

	// Update object metadata to remove the object
	metadataFilePath := filepath.Join(bucketDir, "objects.csv")
	err = removeObjectMetadata(metadataFilePath, objectKey)
	if err != nil {
		// If there is an error updating metadata, return 500 with XML error message
		respondWithError(w, "Error updating object metadata", http.StatusInternalServerError)
		log.Printf("Error updating metadata: %v", err)
		return
	}

	// Respond with success message in XML format
	respondWithSuccess(w, fmt.Sprintf("Object '%s' deleted successfully", objectKey))
}
