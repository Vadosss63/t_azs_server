package application

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (a app) listLogFiles(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		http.Error(rw, "Ivalid token: ", http.StatusInternalServerError)
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
	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	_, ok := getIntVal(id)

	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "error id"})
		return
	}

	// 10 MB максимальный размер файла
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", err.Error()})
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", err.Error()})
		return
	}
	defer file.Close()

	uploadsDir := "./uploads/" + id + "/"
	if _, err := os.Stat(uploadsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadsDir, 0755); err != nil { // Используйте os.MkdirAll вместо os.Mkdir
			rw.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(rw).Encode(answer{"error", err.Error()})
			return
		}
	}

	dst, err := os.Create(uploadsDir + handler.Filename)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(answer{"error", err.Error()})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(answer{"error", err.Error()})
		return
	}

	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(answer{"ok", "Файл успешно загружен"})
}

func (a app) getLogButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		// a.repo.UpdateLogButton(a.ctx, idInt, 1)
		logButton, err := a.repo.GetLogButton(a.ctx, idInt)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(logButton)
			return
		}
	}
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(answer{Msg: "error", Status: "error id or GetLogButton"})
}

func (a app) resetLogButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	a.resetLogAzs(rw, r, p)
}

func (a app) setLogCmd(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Set("Content-Type", "application/json")
	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.UpdateLogButton(a.ctx, idInt, 1)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(answer{Msg: "Ok", Status: "Ok"})
			return
		}
	}
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(answer{Msg: "error", Status: "error"})
}

func (a app) resetLogAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		err := a.repo.UpdateLogButton(a.ctx, idInt, 0)
		if err == nil {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(answer{Msg: "Ok", Status: "Ok"})
			return
		}
	}
	rw.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(rw).Encode(answer{Msg: "error", Status: "error"})
}

func (a app) deleteLogs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	rw.Header().Set("Content-Type", "application/json")

	tokenReq := strings.TrimSpace(r.FormValue("token"))
	if a.token != tokenReq {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "invalid token"})
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(answer{"error", "ID is required"})
		return
	}

	uploadsDir := filepath.Join("./uploads", id) + "/"

	err := os.RemoveAll(uploadsDir)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(rw).Encode(answer{"error", err.Error()})
		return
	}

	http.Redirect(rw, r, "/list_logs", http.StatusSeeOther)
}
