package handler

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"
)

// Helper function to check if a bucket exists in the CSV.
func bucketExists(bucketName string) bool {
	file, err := os.Open(BucketFile)
	if err != nil {
		log.Printf("Error opening metadata file: %v", err)
		return false
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("Error reading metadata file: %v", err)
		return false
	}

	for _, record := range records {
		if record[0] == bucketName {
			return true // Bucket found regardless of status
		}
	}
	return false
}

// Helper function to append a new bucket to the CSV file.
func appendToCSV(bucketName, status string) error {
	file, err := os.OpenFile(BucketFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return fmt.Errorf("error opening metadata file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	creationTime := time.Now().Format(time.RFC3339)
	record := []string{bucketName, creationTime, creationTime, status}

	err = writer.Write(record)
	if err != nil {
		return fmt.Errorf("error writing to metadata file: %w", err)
	}
	return nil
}

// Helper function to update the status of a bucket in the CSV file.
func updateBucketStatus(bucketName, newStatus string) error {
	file, err := os.Open(BucketFile)
	if err != nil {
		return fmt.Errorf("error opening metadata file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("error reading metadata file: %w", err)
	}

	// Update the status for the specific bucket
	for i, record := range records {
		if record[0] == bucketName && record[3] == "Active" {
			record[3] = newStatus
			record[2] = time.Now().Format(time.RFC3339) // Update the last modified time
			records[i] = record
			break
		}
	}

	// Rewrite the CSV file with updated records
	file, err = os.OpenFile(BucketFile, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("error opening metadata file for writing: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(records)
	if err != nil {
		return fmt.Errorf("error writing updated metadata file: %w", err)
	}

	return nil
}
