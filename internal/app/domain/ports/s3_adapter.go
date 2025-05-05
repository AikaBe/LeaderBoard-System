package ports

import "net/http"

type S3Adapter interface {
	UploadImage(r *http.Request, imageType string) (string, error)
}
