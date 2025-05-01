package s3

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type Adapter struct {
	TripleSBaseURL string // например: http://triple-s:9090
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

	// Считываем содержимое файла в память
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return "", err
	}

	// Сформировать URL: images/{photo}
	uploadURL := fmt.Sprintf("%s/images/%s", a.TripleSBaseURL, header.Filename)

	// Создать PUT-запрос с телом = содержимому файла
	req, err := http.NewRequest(http.MethodPut, uploadURL, &buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	// Отправить запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Считать ответ от Triple-S
	respBuf := new(bytes.Buffer)
	respBuf.ReadFrom(resp.Body)

	return respBuf.String(), nil // предполагается, что triple-s возвращает URL как текст
}
