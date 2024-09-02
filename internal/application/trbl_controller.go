package application

import (
	"fmt"
	"html/template"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

const (
	logsPath      = "/tmp/t_azs/"
	maxUploadSize = 10 << 20 // 10 MB
)

type LogsPageTemplate struct {
	IdAzs  string
	IdUser int
	Logs   []string
}

func processFileUpload(r *http.Request, maxUploadSize int64) (multipart.File, *multipart.FileHeader, error) {
	contentType := r.Header.Get("Content-Type")
	log.Println("Received Content-Type:", contentType)

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		log.Printf("Error parsing media type '%s': %s\n", contentType, err)
		return nil, nil, fmt.Errorf("Error parsing media type: %s", err)
	}

	if mediaType != "multipart/form-data" {
		log.Printf("Expected 'multipart/form-data', but got '%s'\n", mediaType)
		return nil, nil, fmt.Errorf("Invalid Content-Type: Expected multipart/form-data, got: %s", mediaType)
	}

	boundary, ok := params["boundary"]
	if !ok {
		log.Println("No boundary parameter found in Content-Type")
		return nil, nil, fmt.Errorf("No boundary found in Content-Type")
	}
	log.Printf("Boundary received: %s\n", boundary)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Printf("Error parsing multipart form with boundary '%s': %s\n", boundary, err)
		return nil, nil, fmt.Errorf("Error parsing multipart form: %s", err)
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("Error retrieving the file: %s\n", err)
		return nil, nil, err
	}

	if handler.Size > maxUploadSize {
		log.Printf("File too large: %d bytes, maximum allowed: %d bytes\n", handler.Size, maxUploadSize)
		file.Close()
		return nil, nil, fmt.Errorf("File too large: %d bytes, maximum allowed: %d bytes", handler.Size, maxUploadSize)
	}

	return file, handler, nil
}

func (a app) uploadLogs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := strings.TrimSpace(r.FormValue("id"))
	_, ok := getIntVal(id)

	if !ok {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	file, handler, err := processFileUpload(r, maxUploadSize)
	if err != nil {
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	defer file.Close()

	uploadsDir := filepath.Join(logsPath, id)

	err = saveUploadedFile(uploadsDir, handler.Filename, file)
	if err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	sendJsonResponse(rw, http.StatusOK, "Файл успешно загружен", "Ok")
}

func (a app) getLogButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	idInt, ok := getIntVal(strings.TrimSpace(r.FormValue("id")))

	if !ok {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}
	logButton, err := a.repo.TrblButtonRepo.GetLogButton(a.ctx, idInt)
	if err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")

		return
	}
	sendJson(rw, http.StatusOK, logButton)
}

func (a app) resetLogButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	a.logButtonReset(rw, r, p)
}

func (a app) downloadLogFile(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fileName := strings.TrimSpace(r.FormValue("file"))
	id := strings.TrimSpace(r.FormValue("id_azs"))
	filePath := logsPath + id + "/" + fileName

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(rw, "Файл не найден", http.StatusNotFound)
		return
	}

	rw.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	rw.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(rw, r, filePath)
}

func (a app) listLogFiles(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	if _, ok := getIntVal(id); !ok {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	uploadsDir := logsPath + id + "/"
	if err := ensureDirectory(uploadsDir); err != nil {
		http.Error(rw, "Не удалось создать директорию: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileNames, err := listFilesInDirectory(uploadsDir)
	if err != nil {
		http.Error(rw, "Не удалось прочитать директорию: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logsPageTemplate := LogsPageTemplate{
		IdAzs:  id,
		IdUser: 0, // Установите корректное значение IdUser, если это необходимо
		Logs:   fileNames,
	}

	tpl, err := template.ParseFiles(
		filepath.Join("public", "html", "logs_page.html"),
		filepath.Join("public", "html", "admin_navi.html"),
	)
	if err != nil {
		http.Error(rw, "Ошибка при парсинге шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(rw, "logsPageTemplate", logsPageTemplate); err != nil {
		http.Error(rw, "Ошибка при рендеринге шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}

func (a app) deleteLogs(rw http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	if id == "" {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	uploadsDir := filepath.Join(logsPath, id) + "/"

	if err := deleteDirectory(uploadsDir); err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	http.Redirect(rw, r, "/list_logs?id_azs="+id, http.StatusSeeOther)
}

func (a app) logButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	pushedBtn := r.FormValue("pushedBtn")

	switch pushedBtn {
	case "download":
		a.setLogCmd(rw, r)
		return
	case "delete":
		a.deleteLogs(rw, r)
		return
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")

}

func (a app) setLogCmd(rw http.ResponseWriter, r *http.Request) {

	id := strings.TrimSpace(r.FormValue("id_azs"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.TrblButtonRepo.UpdateLogButton(a.ctx, idInt, 1)
		if err == nil {
			sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (a app) logButtonReady(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idInt, ok := getIntVal(r.FormValue("id_azs"))
	if !ok {
		sendError(rw, "Invalid id_azs: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	button, err := a.repo.TrblButtonRepo.GetLogButton(a.ctx, idInt)
	if err != nil {
		sendError(rw, "Error fetching update button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if button.Download == 0 {
		sendJsonResponse(rw, http.StatusOK, "Ok", "ready")
	} else {
		sendJsonResponse(rw, http.StatusOK, "Ok", "not_ready")
	}
}

func (a app) logButtonReset(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.TrblButtonRepo.UpdateLogButton(a.ctx, idInt, 0)
		if err == nil {
			sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}
