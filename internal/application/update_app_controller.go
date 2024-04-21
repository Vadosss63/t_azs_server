package application

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"strconv"

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

func (a app) getAppUpdateButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)

	if ok {
		updateCommand, err := a.repo.GetUpdateCommand(a.ctx, idInt)
		if err == nil {
			sendJson(rw, http.StatusOK, updateCommand)
			return
		}
	}
	sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (a app) setAppUpdateCmd(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
    id := strings.TrimSpace(r.FormValue("id_azs"))
    idInt, ok := getIntVal(id)

    version := strings.TrimSpace(r.FormValue("version"))

    if ok {
        url := "http://t-azs.ru:" + strconv.Itoa(a.port) + "/install/" + version
        err := a.repo.UpdateUpdateCommand(a.ctx, idInt, url)
        if err == nil {
            sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")
            return
        }
    }
    sendJsonResponse(rw, http.StatusBadRequest, "Error", "Error")
}

func (a app) resetAppUpdateButton(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !a.validateToken(rw, r.FormValue("token")) {
		return
	}

	a.resetAppUpdateAzs(rw, r, p)
}

func (a app) resetAppUpdateAzs(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	idInt, ok := getIntVal(strings.TrimSpace(r.FormValue("id")))

	if !ok {
		sendJsonResponse(rw, http.StatusBadRequest, "Error id", "Error")
		return
	}

	err := a.repo.UpdateUpdateCommand(a.ctx, idInt, "")
	if err != nil {
		sendJsonResponse(rw, http.StatusInternalServerError, err.Error(), "Error")

		return
	}
	sendJsonResponse(rw, http.StatusOK, "Ok", "Ok")

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

	err = saveUploadedFile(destinationDir, filename, response.Body)
	if err != nil {
		return fmt.Errorf("error saving downloaded file: %v", err)
	}

	return nil
}

func (a app) downloadVersionHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	if _, ok := getIntVal(id); !ok {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}
	version := r.URL.Query().Get("version")
	filename := fmt.Sprintf("%s.tar.gz", version)
	filePath := filepath.Join(updateAppPath, filename)

	exists, err := checkFileExists(filePath)
	if err != nil {
		http.Error(rw, "Ошибка сервера при проверке файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		err = downloadFromGitHub("Vadosss63", "GasStationPro", version, updateAppPath)
		if err != nil {
			http.Error(rw, "Ошибка при скачивании файла: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(rw, r, "/update_app_page?id_azs="+id, http.StatusSeeOther)
}

func (a app) showUpdateAppPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id_azs"))
	if _, ok := getIntVal(id); !ok {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	tags, err := fetchRepositoryTags("Vadosss63", "GasStationPro") // Адаптируйте к вашему репозиторию
	if err != nil {
		http.Error(rw, "Failed to fetch repository tags: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ensureDirectory(updateAppPath); err != nil {
		http.Error(rw, "Не удалось создать директорию: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileNames, err := listFilesInDirectory(updateAppPath)
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
