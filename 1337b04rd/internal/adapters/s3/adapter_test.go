package s3_test

import (
	"1337b04rd/internal/adapters/s3"
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUploadImage_Success(t *testing.T) {
	// Фейковый сервер для эмуляции TripleS
	s3Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			http.NotFound(w, r) // симулируем, что бакет не существует
		case http.MethodPut:
			// Разделяем по URL: создание бакета и загрузка изображения
			if strings.Contains(r.URL.Path, "post-images/test-image.png") {
				// Это загрузка изображения
				w.WriteHeader(http.StatusOK)
			} else {
				// Это создание бакета
				w.WriteHeader(http.StatusCreated)
			}
		default:
			http.Error(w, "unexpected method", http.StatusMethodNotAllowed)
		}
	}))

	publicURL := "http://public-url.com"
	adapter := s3.NewAdapter(s3Server.URL, publicURL)

	// Создание multipart-запроса с изображением
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", "test-image.png")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, strings.NewReader("fake image content"))
	if err != nil {
		t.Fatalf("failed to write image data: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Вызов метода UploadImage
	resultURL, err := adapter.UploadImage(req, "post")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	expectedURL := publicURL + "/post-images/test-image.png"
	if resultURL != expectedURL {
		t.Errorf("unexpected URL: got %s, want %s", resultURL, expectedURL)
	}
}
