package s3

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Adapter struct {
	TripleSBaseURL  string
	PublicAccessURL string
}

// NewAdapter creates a new Adapter instance.
func NewAdapter(internalURL, publicURL string) *Adapter {
	return &Adapter{
		TripleSBaseURL:  internalURL,
		PublicAccessURL: publicURL,
	}
}

// UploadImage uploads an image from the HTTP request to the appropriate bucket.
// Returns the public URL of the uploaded image or an error.
func (a *Adapter) UploadImage(r *http.Request, imageType string) (string, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	// Determine bucket name based on image type
	bucketName := "images"
	switch imageType {
	case "post":
		bucketName = "post-images"
	case "comment":
		bucketName = "images"
	default:
		bucketName = "misc-images"
	}

	// Check if the bucket exists
	checkBucketURL := fmt.Sprintf("%s/%s", a.TripleSBaseURL, bucketName)
	checkReq, err := http.NewRequest(http.MethodGet, checkBucketURL, nil)
	if err != nil {
		slog.Error("Failed to create bucket check request", "error", err)
		return "", err
	}
	checkResp, err := http.DefaultClient.Do(checkReq)
	if err != nil {
		slog.Error("Failed to check bucket existence", "error", err)
		return "", err
	}
	defer checkResp.Body.Close()

	// Create bucket if it doesn't exist
	if checkResp.StatusCode == http.StatusNotFound {
		slog.Info("Bucket does not exist. Creating...", "bucket", bucketName)

		createReq, err := http.NewRequest(http.MethodPut, checkBucketURL, nil)
		if err != nil {
			slog.Error("Failed to create bucket request", "error", err)
			return "", err
		}
		createResp, err := http.DefaultClient.Do(createReq)
		if err != nil {
			slog.Error("Failed to send bucket creation request", "error", err)
			return "", err
		}
		defer createResp.Body.Close()

		if createResp.StatusCode != http.StatusCreated && createResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(createResp.Body)
			slog.Error("Bucket creation failed", "response", string(body))
			return "", fmt.Errorf("bucket creation failed: %s", body)
		}
		slog.Info("Bucket created successfully", "bucket", bucketName)
	}

	// Read the uploaded file into buffer
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, file); err != nil {
		slog.Error("Failed to read file", "error", err)
		return "", err
	}

	// Upload the image to the bucket
	uploadURL := fmt.Sprintf("%s/%s/%s", a.TripleSBaseURL, bucketName, header.Filename)
	publicURL := fmt.Sprintf("%s/%s/%s", a.PublicAccessURL, bucketName, header.Filename)

	uploadReq, err := http.NewRequest(http.MethodPut, uploadURL, &buf)
	if err != nil {
		slog.Error("Failed to create upload request", "error", err)
		return "", err
	}
	uploadReq.Header.Set("Content-Type", "application/octet-stream")

	uploadResp, err := http.DefaultClient.Do(uploadReq)
	if err != nil {
		slog.Error("Failed to upload image", "error", err)
		return "", err
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(uploadResp.Body)
		slog.Error("Image upload failed", "response", string(body))
		return "", fmt.Errorf("upload failed: %s", body)
	}

	slog.Info("Image uploaded successfully", "url", publicURL)
	return publicURL, nil
}
