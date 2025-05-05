package bucketHandler

import (
	"net/http"
	"strings"

	"triple-s/handler"
)

func BucketHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the path segments
	pathSegments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	// Check if we have at least the bucket name
	if len(pathSegments) < 1 {
		http.Error(w, "BucketName is required", http.StatusBadRequest)
		return
	}

	// Get the method (the first segment) and the bucket name (the second segment)
	method := r.Method
	bucketName := pathSegments[0]
	var objectKey string
	if len(pathSegments) > 2 {
		// Reject if there are more than 2 segments
		http.Error(w, "Invalid path. Only 2 path segments allowed.", http.StatusBadRequest)
		return
	}

	// If there's more than 1 segment, the second one is the object key
	if len(pathSegments) > 1 {
		objectKey = pathSegments[1]
	}

	// Switch based on the HTTP method (PUT, GET, DELETE)
	switch method {

	case http.MethodPut:
		// If there's an object key, handle object upload (PUT with objectKey)
		if objectKey != "" {
			handler.UploadObjectHandler(w, r)
		} else {
			// Otherwise, handle bucket creation (PUT with bucketName)
			handler.CreateBucketHandler(w, r, bucketName)
		}

	case http.MethodDelete:
		// If there's an object key, handle object deletion (DELETE with objectKey)
		if objectKey != "" {
			handler.DeleteObjectHandler(w, r)
		} else {
			// Otherwise, handle bucket deletion (DELETE with bucketName)
			handler.DeleteBucketHandler(w, r, bucketName)
		}

	case http.MethodGet:
		// If there's an object key, handle object retrieval (GET with objectKey)
		if objectKey != "" {
			handler.GetObjectHandler(w, r)
		} else {
			// Otherwise, handle listing of buckets (GET with bucketName)
			handler.ListBucketsHandler(w, r)
		}

	default:
		// If the method is not supported, return an error
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
