package application

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (a app) listLogFiles(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	_, ok := getIntVal(id)

	if !ok {
		http.Error(rw, "Ivalid token: ", http.StatusInternalServerError)
		return
	}

	uploadsDir := "./uploads/" + id + "/"

	files, err := os.ReadDir(uploadsDir)
	if err != nil {
		http.Error(rw, "Не удалось прочитать директорию: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(rw, "<h1>Logs:</h1>")
	for _, file := range files {
		fmt.Fprintf(rw, "<a href='/download_log?file=%s&id=%s&token=%s'>%s</a><br>", file.Name(), id, a.token, file.Name())
	}
}

func (a app) downloadLogFile(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fileName := strings.TrimSpace(r.FormValue("file"))
	id := strings.TrimSpace(r.FormValue("id"))
	filePath := "./uploads/" + id + "/" + fileName

	// Проверяем, существует ли файл
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(rw, "Файл не найден", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовки для скачивания файла
	rw.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	rw.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(rw, r, filePath)
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

	// 10 MB максимальный размер файла
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		sendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	defer file.Close()

	uploadsDir := "./uploads/" + id + "/"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadsDir, 0755); err != nil {
			sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
			return
		}
	}

	dst, err := os.Create(uploadsDir + handler.Filename)
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

func (a app) setLogCmd(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}
	id := strings.TrimSpace(r.FormValue("id"))
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

func (a app) deleteLogs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	uploadsDir := filepath.Join("./uploads", id) + "/"

	err := os.RemoveAll(uploadsDir)
	if err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	http.Redirect(rw, r, "/list_logs", http.StatusSeeOther)
}
