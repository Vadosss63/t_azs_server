package admin_controller

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/Vadosss63/t-azs/internal/application"
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

type AdminController struct {
	app *application.App
}

func NewController(app *application.App) *AdminController {
	return &AdminController{app: app}
}

func (c AdminController) Routes(router *httprouter.Router) {
	router.POST("/add_user_to_asz", c.app.Authorized(c.addUserToAsz))

	router.GET("/users", c.app.Authorized(c.showUsersPage))

	router.POST("/show_azs_for", c.app.Authorized(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id_user, ok_id := application.GetIntVal(r.FormValue("user"))

		userId, ok := r.Context().Value("userId").(int)

		if !ok || !ok_id {
			http.Error(rw, "Error user", http.StatusBadRequest)
			return
		}
		u, err := c.app.Repo.UserRepo.Get(c.app.Ctx, userId)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		c.AdminPage(rw, r, p, u, id_user)
	}))
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

func (c AdminController) AdminPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, u user.User, id int) {

	// c.app.Repo.CreateYaAzsInfoTable(c.app.Ctx)

	var azs_statses []azs.AzsStatsData
	var err error

	if id == -2 {
		azs_statses, err = c.app.Repo.AzsRepo.GetAll(c.app.Ctx)
	} else {
		azs_statses, err = c.app.Repo.AzsRepo.GetAzsAllForUser(c.app.Ctx, id)
	}

	users, err := c.app.Repo.UserRepo.GetAll(c.app.Ctx)
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
		azsStatsDataFull.IsEnabled, err = c.app.Repo.YaAzsRepo.GetEnable(c.app.Ctx, azsStatsDataFull.IdAzs)

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

func (c AdminController) showUsersPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	users, err := c.app.Repo.UserRepo.GetAll(c.app.Ctx)
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

func (c AdminController) addUserToAsz(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id_azs, _ := application.GetIntVal(r.FormValue("id_azs"))
	id_user, _ := application.GetIntVal(r.FormValue("user"))

	err := c.app.Repo.AzsRepo.AddAzsToUser(c.app.Ctx, id_user, id_azs)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}
