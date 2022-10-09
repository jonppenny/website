package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"github.com/justinas/nosurf"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	err = app.errorLog.Output(2, trace)
	if err != nil {
		log.Fatal(err)
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData, isAdmin, isCredentials bool) {
	buf := new(bytes.Buffer)

	switch {
	case isAdmin:
		adminTemplates, ok := app.adminTemplateCache[name]
		if !ok {
			app.serverError(w, fmt.Errorf("the template %s does not exist in admin", name))
			return
		}

		err := adminTemplates.Execute(buf, app.addDefaultData(td, r))
		if err != nil {
			app.serverError(w, err)
		}
		break
	case isCredentials:
		credentialsTemplates, ok := app.credentialsTemplateCache[name]
		if !ok {
			app.serverError(w, fmt.Errorf("the template %s does not exist in credentials", name))
			return
		}

		err := credentialsTemplates.Execute(buf, app.addDefaultData(td, r))
		if err != nil {
			app.serverError(w, err)
		}
		break
	default:
		websiteTemplates, ok := app.templateCache[name]
		if !ok {
			app.serverError(w, fmt.Errorf("the template %s does not exist", name))
			return
		}

		err := websiteTemplates.Execute(buf, app.addDefaultData(td, r))
		if err != nil {
			app.serverError(w, err)
		}
		break
	}

	_, err := buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(contextKeyIsAuthenticated).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) uploadFile(r *http.Request, fileInput string) (string, error) {
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

	err = os.MkdirAll("web/static/media/", 0755)
	if err != nil {
		return "", err
	}

	// Create file
	dst, err := os.Create(filepath.Join("web/static/media/", filepath.Base(handler.Filename)))
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

func (app *application) readInt(qs url.Values, key string, def int) (int, error) {
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
