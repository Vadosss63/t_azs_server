package application

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/Vadosss63/t-azs/internal/repository/azs"
	"github.com/Vadosss63/t-azs/internal/repository/user"
	"github.com/julienschmidt/httprouter"
)

type AdminPageTemplate struct {
	User           user.User
	Users          []user.User
	Azses          []azs.AzsStatsDataFull
	SelectedUserId int
}

func deleteUser(users []user.User, login string) []user.User {
	for i, user := range users {
		if user.Login == login {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	return users
}

func (a app) adminPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, u user.User, id int) {

	// a.repo.CreateYaAzsInfoTable(a.ctx)

	var azs_statses []azs.AzsStatsData
	var err error

	if id == -2 {
		azs_statses, err = a.repo.AzsRepo.GetAzsAll(a.ctx)
	} else {
		azs_statses, err = a.repo.AzsRepo.GetAzsAllForUser(a.ctx, id)
	}

	users, err := a.repo.UserRepo.GetUserAll(a.ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	users = deleteUser(users, "admin")

	adminPageTemplate := AdminPageTemplate{
		User:           u,
		Users:          users,
		Azses:          []azs.AzsStatsDataFull{},
		SelectedUserId: id,
	}

	for _, azs_stats := range azs_statses {
		azsStatsDataFull, err := azs.ParseStats(azs_stats)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		azsStatsDataFull.IsEnabled, err = a.repo.GetYaAzsInfoEnable(a.ctx, azsStatsDataFull.IdAzs)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		adminPageTemplate.Azses = append(adminPageTemplate.Azses, azsStatsDataFull)
	}

	lp := filepath.Join("public", "html", "admin_page.html")
	navi := filepath.Join("public", "html", "admin_navi.html")
	azs := filepath.Join("public", "html", "azs_container.html")
	tmpl := template.Must(template.ParseFiles(lp, navi, azs))

	err = tmpl.ExecuteTemplate(rw, "AdminPageTemplate", adminPageTemplate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) showUsersPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	users, err := a.repo.UserRepo.GetUserAll(a.ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "users_page.html")
	navi := filepath.Join("public", "html", "admin_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))
	err = tmpl.ExecuteTemplate(rw, "User", users)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) addUserToAsz(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id_azs, _ := getIntVal(r.FormValue("id_azs"))
	id_user, _ := getIntVal(r.FormValue("user"))

	err := a.repo.AzsRepo.AddAzsToUser(a.ctx, id_user, id_azs)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (a app) showUsersAzsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	userId, ok := getIntVal(r.FormValue("user"))

	if !ok {
		http.Error(rw, "Error userId", http.StatusBadRequest)
		return
	}

	u, err := a.repo.UserRepo.GetUser(a.ctx, userId)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	a.userPage(rw, r, p, u)
}
