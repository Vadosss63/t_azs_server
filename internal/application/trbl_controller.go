package application

import (
	"html/template"
	"io"
	"log"
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

func ensureDirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (a app) uploadLogs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	_, ok := getIntVal(id)
	if !ok {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	contentType := r.Header.Get("Content-Type")
	log.Println("Content-Type:", contentType)
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		log.Println("Invalid Content-Type: Expected multipart/form-data")
		sendJsonResponse(rw, http.StatusBadRequest, "Invalid Content-Type: Expected multipart/form-data", "Error")
		return
	}

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		log.Println(err.Error())
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	defer file.Close()

	if handler.Size > maxUploadSize {
		sendJsonResponse(rw, http.StatusBadRequest, "Файл слишком большой", "Error")
		return
	}

	uploadsDir := filepath.Join(logsPath, id)
	if err := ensureDirectory(uploadsDir); err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	safeFilename := filepath.Base(handler.Filename)
	dst, err := os.Create(filepath.Join(uploadsDir, safeFilename))
	if err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	sendJsonResponse(rw, http.StatusOK, "Файл успешно загружен", "Ok")
}

func (a app) getLogButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}

	idInt, ok := getIntVal(strings.TrimSpace(r.FormValue("id")))

	if !ok {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}
	logButton, err := a.repo.GetLogButton(a.ctx, idInt)
	if err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")

		return
	}
	sendJson(rw, http.StatusOK, logButton)
}

func (a app) resetLogButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}

	a.resetLogAzs(rw, r, p)
}

func (a app) resetLogAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.UpdateLogButton(a.ctx, idInt, 0)
		if err == nil {
			sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (a app) setLogCmd(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.UpdateLogButton(a.ctx, idInt, 1)
		if err == nil {
			sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
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

	files, err := os.ReadDir(uploadsDir)
	if err != nil {
		http.Error(rw, "Не удалось прочитать директорию: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logsPageTemplate := LogsPageTemplate{
		IdAzs:  id,
		IdUser: 0, // Установите корректное значение IdUser, если это необходимо
		Logs:   make([]string, 0, len(files)),
	}
	for _, file := range files {
		if !file.IsDir() { // Убедитесь, что это файл, а не директория
			logsPageTemplate.Logs = append(logsPageTemplate.Logs, file.Name())
		}
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

func (a app) deleteLogs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	if id == "" {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	uploadsDir := filepath.Join(logsPath, id) + "/"

	err := os.RemoveAll(uploadsDir)
	if err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	http.Redirect(rw, r, "/list_logs?id_azs="+id, http.StatusSeeOther)
}
