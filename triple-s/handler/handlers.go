package handler

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"triple-s/config"
)

// Bucket represents the metadata of a bucket.
type Bucket struct {
	Name             string `xml:"Name"`
	CreationTime     string `xml:"CreationTime"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"`
}

// const metadataFilePath = "data/buckets.csv"

var (
	StorageDir string
	BucketFile string
)

// InitStorage initializes the storage path and bucket file path
func InitStorage() {
	// Initialize the global variables using the global config
	StorageDir = config.GetStorage() // Get the storage directory from the global config
	BucketFile = filepath.Join(StorageDir, "buckets.csv")

	// Create the storage directory and CSV file if they don't exist
	err := CreateStorageDirAndCSV(StorageDir, BucketFile)
	if err != nil {
		log.Fatalf("Error initializing storage: %v", err)
	}

	// Log the paths for debugging purposes
	log.Printf("Storage directory: %s", StorageDir)
	log.Printf("Bucket file: %s", BucketFile)
}

func ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the bucket name from the URL path (if present)
	bucketName := strings.TrimPrefix(r.URL.Path, "/")

	// Open the CSV file
	file, err := os.Open(BucketFile)
	if err != nil {
		http.Error(w, "Error opening metadata storage", http.StatusInternalServerError)
		log.Printf("Error opening metadata storage: %v", err)
		return
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Error reading metadata storage", http.StatusInternalServerError)
		log.Printf("Error reading metadata storage: %v", err)
		return
	}

	// Check if we need to list all buckets or a specific bucket
	if bucketName == "" {
		// List all buckets
		var buckets []Bucket
		for _, record := range records {
			// Assuming record format: Name, CreationTime, LastModifiedTime, Status
			bucket := Bucket{
				Name:             record[0],
				CreationTime:     record[1],
				LastModifiedTime: record[2],
				Status:           record[3],
			}
			buckets = append(buckets, bucket)
		}

		// Wrap the buckets in a response structure
		type BucketList struct {
			XMLName xml.Name `xml:"BucketList"`
			Buckets []Bucket `xml:"Bucket"`
		}

		bucketList := BucketList{Buckets: buckets}

		// Respond with the XML representation of all buckets
		w.Header().Set("Content-Type", "application/xml")
		xmlBytes, err := xml.MarshalIndent(bucketList, "", "  ")
		if err != nil {
			http.Error(w, "Error marshaling bucket data to XML", http.StatusInternalServerError)
			log.Printf("Error marshaling bucket data to XML: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(xmlBytes)
	} else {
		// Get a specific bucket's info
		var bucket Bucket
		var found bool
		for _, record := range records {
			if record[0] == bucketName {
				bucket = Bucket{
					Name:             record[0],
					CreationTime:     record[1],
					LastModifiedTime: record[2],
					Status:           record[3],
				}
				found = true
				break
			}
		}

		if !found {
			http.Error(w, "Bucket not found", http.StatusNotFound)
			return
		}

		// Respond with the XML representation of the specific bucket
		w.Header().Set("Content-Type", "application/xml")
		xmlBytes, err := xml.MarshalIndent(bucket, "", "  ")
		if err != nil {
			http.Error(w, "Error marshaling bucket data to XML", http.StatusInternalServerError)
			log.Printf("Error marshaling bucket data to XML: %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(xmlBytes)
	}
}

// CreateBucketHandler handles bucket creation and updates the CSV file.
func CreateBucketHandler(w http.ResponseWriter, r *http.Request, bucketName string) {
	// Validate bucket name
	if !isValidBucketName(bucketName) {
		http.Error(w, "Invalid bucket name", http.StatusBadRequest)
		return
	}

	// Check if the bucket already exists
	if bucketExists(bucketName) {
		http.Error(w, "Bucket already exists", http.StatusConflict)
		return
	}

	// Create the bucket folder (representing the bucket)
	err := os.Mkdir(filepath.Join(StorageDir, bucketName), 0o755)
	if err != nil {
		http.Error(w, "Error creating bucket folder", http.StatusInternalServerError)
		log.Printf("Error creating bucket folder: %v", err)
		return
	}

	// Create the CSV file for metadata
	err = appendToCSV(bucketName, "Active")
	if err != nil {
		http.Error(w, "Error creating metadata CSV", http.StatusInternalServerError)
		log.Printf("Error creating metadata CSV: %v", err)
		return
	}
	type CreateResponse struct {
		XMLName      xml.Name `xml:"CreateBucketResponse"`
		BucketName   string   `xml:"BucketName"`
		CreationTime string   `xml:"CreationTime"`
		LastModified string   `xml:"LastModifiesTime"`
		Status       string   `xml:"Status"`
	}

	creationTime := time.Now().Format(time.RFC3339)
	bucketResponse := CreateResponse{
		BucketName:   bucketName,
		CreationTime: creationTime,
		LastModified: creationTime, // Initially, creation time and last modified time are the same
		Status:       "Active",
	}

	xmlBytes, err := xml.MarshalIndent(bucketResponse, "", " ")
	if err != nil {
		http.Error(w, "Error marshaling bucket data to XML", http.StatusInternalServerError)
		log.Printf("Error marshaling bucket data to XML: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(xmlBytes)
	if err != nil {
		http.Error(w, "Error writing XML response", http.StatusInternalServerError)
		log.Printf("Error writing XML response: %v", err)
		return
	}
}

func DeleteBucketHandler(w http.ResponseWriter, r *http.Request, bucketName string) {
	// Check if the bucket exists
	if !bucketExists(bucketName) {
		// Respond with an error if the bucket doesn't exist
		http.Error(w, "Bucket not found", http.StatusNotFound)
		return
	}

	// Define the bucket path
	bucketPath := filepath.Join(StorageDir, bucketName)

	// Ensure the bucket is empty before deletion
	if !isBucketEmpty(bucketPath) {
		// If the bucket is not empty, prevent deletion and return an error
		http.Error(w, "Bucket is not empty. Please delete all objects before deleting the bucket.", http.StatusConflict)
		return
	}

	// Optionally delete metadata file explicitly before deleting the folder
	err := deleteBucketMetadata(bucketName)
	if err != nil {
		http.Error(w, "Error deleting metadata CSV file", http.StatusInternalServerError)
		log.Printf("Error deleting metadata CSV file: %v", err)
		return
	}

	// Delete the bucket folder and all its contents (objects and metadata)
	err = os.RemoveAll(bucketPath)
	if err != nil {
		// If there's an error while deleting the directory, send a failure response
		http.Error(w, "Error deleting bucket", http.StatusInternalServerError)
		log.Printf("Error deleting bucket folder: %v", err)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bucket '%s' deleted successfully", bucketName)
}

// Helper function to check if a bucket is empty
func isBucketEmpty(bucketPath string) bool {
	// Open the directory
	dir, err := os.Open(bucketPath)
	if err != nil {
		log.Printf("Error opening directory %s: %v", bucketPath, err)
		return false
	}
	defer dir.Close()

	// Read the directory contents
	entries, err := dir.Readdirnames(0) // 0 means to read all names
	if err != nil {
		log.Printf("Error reading directory %s: %v", bucketPath, err)
		return false
	}

	// If there are any entries, the bucket is not empty
	return len(entries) == 1
}

// Helper function to delete the bucket's metadata from the CSV file.
func deleteBucketMetadata(bucketName string) error {
	// Open the CSV file for reading and writing
	file, err := os.OpenFile(BucketFile, os.O_RDWR, 0o644) // Open with read and write permissions
	if err != nil {
		return fmt.Errorf("error opening metadata file: %w", err)
	}
	defer file.Close()

	// Read all records from the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading metadata file: %w", err)
	}

	// Create a new slice of records that excludes the bucket to be deleted
	var updatedRecords [][]string
	var header []string

	// Check if there are any records, including the header row
	if len(records) > 0 {
		header = records[0] // Assuming the first record is the header
		// Add the header to the updatedRecords slice (we want to keep the header)
		updatedRecords = append(updatedRecords, header)
	}

	// Loop through the remaining records and remove the bucket to be deleted
	for _, record := range records[1:] { // Skip the first header row
		if record[0] != bucketName { // Skip the bucket with the matching name
			updatedRecords = append(updatedRecords, record)
		}
	}

	// If the bucket was not found in the metadata, return an error
	if len(updatedRecords) == len(records) {
		return fmt.Errorf("bucket '%s' not found in metadata", bucketName)
	}

	// Now rewrite the CSV file with the updated records
	// First, truncate the file and seek to the beginning
	file.Truncate(0) // Clear the contents of the file
	file.Seek(0, 0)  // Rewind the file pointer to the beginning

	// Write the updated records back to the file
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write all updated records back to the CSV file
	err = writer.WriteAll(updatedRecords)
	if err != nil {
		return fmt.Errorf("error writing updated records to CSV: %w", err)
	}

	return nil
}

// CreateStorageDirAndCSV creates the storage directory and the CSV file inside it.
func CreateStorageDirAndCSV(storageDir, bucketFile string) error {
	// Check if the directory exists
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		// If it doesn't exist, create the directory
		err := os.MkdirAll(storageDir, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create storage directory: %w", err)
		}

	}

	// Check if the CSV file exists
	if _, err := os.Stat(bucketFile); os.IsNotExist(err) {
		// If it doesn't exist, create the CSV file
		file, err := os.Create(bucketFile)
		if err != nil {
			return fmt.Errorf("failed to create CSV file: %w", err)
		}
		defer file.Close()

		// Create a CSV writer and write a header (optional)
		writer := csv.NewWriter(file)
		defer writer.Flush()
		header := []string{"Name", "CreationTime", "LastModifiedTime", "Status"}
		err = writer.Write(header)
		if err != nil {
			return fmt.Errorf("failed to write header to CSV: %w", err)
		}

	}

	return nil
}
