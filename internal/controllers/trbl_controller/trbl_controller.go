package trbl_controller

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

	"github.com/Vadosss63/t-azs/internal/application"
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

type TrblControllerController struct {
	app *application.App
}

func NewController(app *application.App) *TrblControllerController {
	return &TrblControllerController{app: app}
}

func (c TrblControllerController) Routes(router *httprouter.Router) {

	router.POST("/get_log_cmd", c.app.Authorized(c.getLogButton))
	router.POST("/upload_log", c.app.Authorized(c.uploadLogs))
	router.POST("/reset_log_cmd", c.app.Authorized(c.resetLogButton))

	router.POST("/log_button", c.app.Authorized(c.logButton))
	router.GET("/log_button_ready", c.app.Authorized(c.logButtonReady))
	router.GET("/log_button_reset", c.app.Authorized(c.logButtonReset))

	router.GET("/list_logs", c.app.Authorized(c.listLogFiles))
	router.GET("/download_log", c.app.Authorized(c.downloadLogFile))

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

func (c TrblControllerController) uploadLogs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := strings.TrimSpace(r.FormValue("id"))
	_, ok := application.GetIntVal(id)

	if !ok {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	file, handler, err := processFileUpload(r, maxUploadSize)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, err.Error(), "Error")
		return
	}
	defer file.Close()

	uploadsDir := filepath.Join(logsPath, id)

	err = application.SaveUploadedFile(uploadsDir, handler.Filename, file)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	application.SendJsonResponse(rw, http.StatusOK, "Файл успешно загружен", "Ok")
}

func (c TrblControllerController) getLogButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	idInt, ok := application.GetIntVal(strings.TrimSpace(r.FormValue("id")))

	if !ok {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}
	logButton, err := c.app.Repo.TrblButtonRepo.Get(c.app.Ctx, idInt)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")

		return
	}
	application.SendJson(rw, http.StatusOK, logButton)
}

func (c TrblControllerController) resetLogButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c.logButtonReset(rw, r, p)
}

func (c TrblControllerController) downloadLogFile(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func (c TrblControllerController) listLogFiles(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	if _, ok := application.GetIntVal(id); !ok {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	uploadsDir := logsPath + id + "/"
	if err := application.EnsureDirectory(uploadsDir); err != nil {
		http.Error(rw, "Не удалось создать директорию: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileNames, err := application.ListFilesInDirectory(uploadsDir)
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

func (c TrblControllerController) deleteLogs(rw http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	if id == "" {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	uploadsDir := filepath.Join(logsPath, id) + "/"

	if err := application.DeleteDirectory(uploadsDir); err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")
		return
	}

	http.Redirect(rw, r, "/list_logs?id_azs="+id, http.StatusSeeOther)
}

func (c TrblControllerController) logButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	pushedBtn := r.FormValue("pushedBtn")

	switch pushedBtn {
	case "download":
		c.setLogCmd(rw, r)
		return
	case "delete":
		c.deleteLogs(rw, r)
		return
	}
	application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")

}

func (c TrblControllerController) setLogCmd(rw http.ResponseWriter, r *http.Request) {

	id := strings.TrimSpace(r.FormValue("id_azs"))
	idInt, ok := application.GetIntVal(id)

	if ok {
		err := c.app.Repo.TrblButtonRepo.Update(c.app.Ctx, idInt, 1)
		if err == nil {
			application.SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (c TrblControllerController) logButtonReady(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idInt, ok := application.GetIntVal(r.FormValue("id_azs"))
	if !ok {
		application.SendError(rw, "Invalid id_azs: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	button, err := c.app.Repo.TrblButtonRepo.Get(c.app.Ctx, idInt)
	if err != nil {
		application.SendError(rw, "Error fetching update button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if button.Download == 0 {
		application.SendJsonResponse(rw, http.StatusOK, "Ok", "ready")
	} else {
		application.SendJsonResponse(rw, http.StatusOK, "Ok", "not_ready")
	}
}

func (c TrblControllerController) logButtonReset(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := application.GetIntVal(id)

	if ok {
		err := c.app.Repo.TrblButtonRepo.Update(c.app.Ctx, idInt, 0)
		if err == nil {
			application.SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
			return
		}
	}
	application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}
