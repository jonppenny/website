package helpers

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
)

func UploadFile(r *http.Request, fileInput string) (string, error) {
	// Maximum upload of 10 MB files
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", err
	}

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile(fileInput)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileMimetype, err := mimetype.DetectFile(handler.Filename)
	if err != nil {
		return "", err
	}

	allowedMimes := []string{"jpg", "jpeg", "png", "gif"}
	fileCheck := StringInSlice(fileMimetype.String(), allowedMimes)
	if fileCheck != true {
		return "", errors.New("filetype is not allowed")
	}

	err = os.MkdirAll("static/media/", 0755)
	if err != nil {
		return "", err
	}

	// Create file
	dst, err := os.Create(filepath.Join("static/media/", filepath.Base(handler.Filename)))
	defer dst.Close()
	if err != nil {
		return "", err
	}

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return handler.Filename, nil
}

func ReadInt(qs url.Values, key string, def int) (int, error) {
	s := qs.Get(key)

	if s == "" {
		return def, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return def, err
	}

	return i, nil
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
