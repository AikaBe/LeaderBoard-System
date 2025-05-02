package s3

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Adapter struct {
	TripleSBaseURL string
}

func NewAdapter(baseURL string) *Adapter {
	return &Adapter{TripleSBaseURL: baseURL}
}

func (a *Adapter) UploadImage(r *http.Request) (string, error) {
	file, header, err := r.FormFile("image")
	if err != nil {
		return "", err
	}
	defer file.Close()

	bucketName := "images"

	// Проверка: существует ли бакет с помощью GET
	checkBucketURL := fmt.Sprintf("%s/%s", a.TripleSBaseURL, bucketName)
	checkReq, err := http.NewRequest(http.MethodGet, checkBucketURL, nil)
	if err != nil {
		log.Printf("failed to create bucket check request:%v", err)
		return "", fmt.Errorf("failed to create bucket check request: %w", err)
	}

	checkResp, err := http.DefaultClient.Do(checkReq)
	if err != nil {
		log.Printf("failed to check bucket existence:%v", err)
		return "", fmt.Errorf("failed to check bucket existence: %w", err)
	}
	checkResp.Body.Close()

	// Если бакета нет (GET вернул 404) — создаем
	if checkResp.StatusCode == http.StatusNotFound {
		createReq, err := http.NewRequest(http.MethodPost, checkBucketURL, nil)
		if err != nil {
			log.Printf("failed to create bucket request:%v", err)
			return "", fmt.Errorf("failed to create bucket request: %w", err)
		}

		createResp, err := http.DefaultClient.Do(createReq)
		if err != nil {
			log.Printf("failed to send bucket creation request:%v", err)
			return "", fmt.Errorf("failed to send bucket creation request: %w", err)
		}
		defer createResp.Body.Close()

		if createResp.StatusCode != http.StatusCreated && createResp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(createResp.Body)
			log.Printf("bucket creation failed: %s", body)
			return "", fmt.Errorf("bucket creation failed: %s", body)
		}
	}

	// Считываем содержимое изображения в память
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		log.Printf("failed to read file:%v", err)
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Загружаем изображение
	uploadURL := fmt.Sprintf("%s/%s/%s", a.TripleSBaseURL, bucketName, header.Filename)
	uploadReq, err := http.NewRequest(http.MethodPut, uploadURL, &buf)
	if err != nil {
		log.Printf("failed to create upload request:%v", err)
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}
	uploadReq.Header.Set("Content-Type", "application/octet-stream")

	uploadResp, err := http.DefaultClient.Do(uploadReq)
	if err != nil {
		log.Printf("failed to upload image:%v", err)
		return "", fmt.Errorf("failed to upload image: %w", err)
	}
	defer uploadResp.Body.Close()

	if uploadResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(uploadResp.Body)
		log.Printf("upload failed: %s", body)
		return "", fmt.Errorf("upload failed: %s", body)
	}

	respBuf := new(bytes.Buffer)
	respBuf.ReadFrom(uploadResp.Body)
	log.Printf("image url: %s", uploadURL)

	return respBuf.String(), nil
}
