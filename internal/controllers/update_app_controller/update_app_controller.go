package update_app_controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/julienschmidt/httprouter"
)

const (
	updateAppPath = "/tmp/t_azs/update/"
)

type UpdatePageTemplate struct {
	IdAzs             string
	IdUser            int
	AvailableAppFiles []string
	AvailableTags     []string
}

type UpdateAppController struct {
	app *application.App
}

func NewController(app *application.App) *UpdateAppController {
	return &UpdateAppController{app: app}
}

func (c UpdateAppController) Routes(router *httprouter.Router) {

	router.POST("/get_app_update_button", c.app.Authorized(c.getAppUpdateButton))
	router.POST("/reset_app_update_button", c.app.Authorized(c.resetAppUpdateButton))
	router.POST("/app_update_button", c.app.Authorized(c.appUpdateButton))
	router.GET("/app_update_button_ready", c.app.Authorized(c.appUpdateButtonReady))
	router.GET("/app_update_button_reset", c.app.Authorized(c.resetAppUpdateAzs))
	router.GET("/update_app_page", c.app.Authorized(c.showUpdateAppPage))
}

func (c UpdateAppController) appUpdateButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	pushedBtn := r.FormValue("pushedBtn")

	switch pushedBtn {
	case "install":
		c.setAppUpdateCmd(rw, r)
		return
	case "download":
		c.downloadVersionHandler(rw, r)
		return
	case "delete":
		c.deleteAppFile(rw, r)
		return
	}
	application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (c UpdateAppController) deleteAppFile(rw http.ResponseWriter, r *http.Request) {
	_, ok := application.GetIntVal(r.FormValue("id_azs"))
	if !ok {
		application.SendError(rw, "Invalid id_azs: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	filename := strings.TrimSpace(r.FormValue("value"))

	if filename == "" {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
		return
	}
	filePath := filepath.Join(updateAppPath, filename)

	exists, err := application.CheckFileExists(filePath)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, "Error", "Ошибка сервера при проверке файла: "+err.Error())
		return
	}

	if exists {
		err = application.DeleteDirectory(filePath)
		if err != nil {
			application.SendJsonResponse(rw, http.StatusInternalServerError, "Error", "Ошибка при удалени файла: "+err.Error())
			return
		}
	}
	application.SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
}

func (c UpdateAppController) setAppUpdateCmd(rw http.ResponseWriter, r *http.Request) {
	idInt, ok := application.GetIntVal(r.FormValue("id_azs"))
	if !ok {
		application.SendError(rw, "Invalid id_azs: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	version := strings.TrimSpace(r.FormValue("value"))

	if version == "" {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
		return
	}

	url := "http://t-azs.ru:" + strconv.Itoa(c.app.Port) + "/install/" + version
	err := c.app.Repo.UpdaterButtonRepo.Update(c.app.Ctx, idInt, url)

	if err != nil {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
		return
	}
	application.SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
}

func (c UpdateAppController) downloadVersionHandler(rw http.ResponseWriter, r *http.Request) {
	_, ok := application.GetIntVal(r.FormValue("id_azs"))
	if !ok {
		application.SendError(rw, "Invalid id_azs: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	version := strings.TrimSpace(r.FormValue("value"))

	if version == "" {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
		return
	}

	filename := fmt.Sprintf("%s.tar.gz", version)
	filePath := filepath.Join(updateAppPath, filename)

	exists, err := application.CheckFileExists(filePath)
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, "Error", "Ошибка сервера при проверке файла: "+err.Error())
		return
	}

	if !exists {
		err = downloadFromGitHub("Vadosss63", "GasStationPro", version, updateAppPath)
		if err != nil {
			application.SendJsonResponse(rw, http.StatusInternalServerError, "Error", "Ошибка при скачивании файла: "+err.Error())
			return
		}
	}
	application.SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
}

func (c UpdateAppController) appUpdateButtonReady(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idInt, ok := application.GetIntVal(r.FormValue("id_azs"))
	if !ok {
		application.SendError(rw, "Invalid id_azs: "+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	updateCommand, err := c.app.Repo.UpdaterButtonRepo.Get(c.app.Ctx, idInt)
	if err != nil {
		application.SendError(rw, "Error fetching update button: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if updateCommand.Url == "" {
		application.SendJsonResponse(rw, http.StatusOK, "Ok", "ready")
	} else {
		application.SendJsonResponse(rw, http.StatusOK, "Ok", "not_ready")
	}
}

func (c UpdateAppController) getAppUpdateButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := application.GetIntVal(id)

	if ok {
		updateCommand, err := c.app.Repo.UpdaterButtonRepo.Get(c.app.Ctx, idInt)
		if err == nil {
			application.SendJson(rw, http.StatusOK, updateCommand)
			return
		}
	}
	application.SendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (c UpdateAppController) resetAppUpdateButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c.resetAppUpdateAzs(rw, r, p)
}

func (c UpdateAppController) resetAppUpdateAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idInt, ok := application.GetIntVal(strings.TrimSpace(r.FormValue("id")))

	if !ok {
		application.SendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	err := c.app.Repo.UpdaterButtonRepo.Update(c.app.Ctx, idInt, "")
	if err != nil {
		application.SendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")

		return
	}
	application.SendJsonResponse(rw, http.StatusOK, "Ok", "Ok")

}

func fetchRepositoryTags(owner, repo string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags", owner, repo)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch tags, status code: %d", response.StatusCode)
	}

	var tags []struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(response.Body).Decode(&tags)
	if err != nil {
		return nil, err
	}

	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}
	return tagNames, nil
}

func downloadFromGitHub(owner, repo, tag, destinationDir string) error {
	archiveUrl := fmt.Sprintf("https://github.com/%s/%s/archive/refs/tags/%s.tar.gz", owner, repo, tag)

	response, err := http.Get(archiveUrl)
	if err != nil {
		return fmt.Errorf("error making request to GitHub: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status from GitHub: %v", response.Status)
	}

	filename := fmt.Sprintf("%s.tar.gz", tag)

	err = application.SaveUploadedFile(destinationDir, filename, response.Body)
	if err != nil {
		return fmt.Errorf("error saving downloaded file: %v", err)
	}

	return nil
}

func (c UpdateAppController) showUpdateAppPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	if _, ok := application.GetIntVal(id); !ok {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	tags, err := fetchRepositoryTags("Vadosss63", "GasStationPro") // Адаптируйте к вашему репозиторию
	if err != nil {
		http.Error(rw, "Failed to fetch repository tags: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := application.EnsureDirectory(updateAppPath); err != nil {
		http.Error(rw, "Не удалось создать директорию: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileNames, err := application.ListFilesInDirectory(updateAppPath)
	if err != nil {
		http.Error(rw, "Не удалось прочитать директорию: "+err.Error(), http.StatusInternalServerError)
		return
	}

	updatePageTemplate := UpdatePageTemplate{
		IdAzs:             id,
		IdUser:            0,
		AvailableAppFiles: fileNames,
		AvailableTags:     tags,
	}

	tpl, err := template.ParseFiles(
		filepath.Join("public", "html", "update_app_page.html"),
		filepath.Join("public", "html", "admin_navi.html"),
	)
	if err != nil {
		http.Error(rw, "Ошибка при парсинге шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.ExecuteTemplate(rw, "updatePageTemplate", updatePageTemplate); err != nil {
		http.Error(rw, "Ошибка при рендеринге шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}
