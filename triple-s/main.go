// main.go
package main

import (
	"flag"
	"log"
	"net/http"

	"triple-s/bucketHandler"
	"triple-s/handler"

	"triple-s/config" // Adjust the import path according to your project structure
)

func main() {
	// Define flags in main.go
	port := flag.String("port", ":9000", "Port to run the server on (default: :9000)")
	storageDir := flag.String("dir", "./storage", "Directory to store data (default: ./storage)")

	// Parse the flags once
	flag.Parse()

	// Load configuration
	config.LoadConfig(*port, *storageDir)
	handler.InitStorage()

	// Handle routes
	http.HandleFunc("/", bucketHandler.BucketHandler)

	// Start the server using the provided port
	log.Fatal(http.ListenAndServe(config.GlobalConfig.Port, nil))
}
