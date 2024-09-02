package user_controller

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/Vadosss63/t-azs/internal/application"
	"github.com/julienschmidt/httprouter"
)

type Answer struct {
	Message string
}

type UserController struct {
	app *application.App
}

func NewController(app *application.App) *UserController {
	return &UserController{app: app}
}

func (c UserController) Routes(router *httprouter.Router) {
	router.GET("/login", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c.loginPage(rw, "")
	})

	router.POST("/login", c.login)

	router.GET("/logout", c.logout)

	router.GET("/signup", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c.signupPage(rw, "")
	})
	router.POST("/signup", c.signup)

	router.DELETE("/user", c.app.Authorized(c.deleteUser))

	router.POST("/reset_password", c.app.Authorized(c.resetPasswordUser))

}

func (c UserController) resetPasswordUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id, ok_id := application.GetIntVal(strings.TrimSpace(r.FormValue("userId")))
	password := strings.TrimSpace(r.FormValue("password"))
	password2 := strings.TrimSpace(r.FormValue("password2"))

	if !ok_id || password == "" || password2 == "" {
		http.Error(rw, "Ошибка обновления пароля пользователя", http.StatusBadRequest)
		return
	}

	if password != password2 {
		http.Error(rw, "Пароли не совпадают! Попробуйте еще", http.StatusBadRequest)
		return
	}

	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])

	err := c.app.Repo.UserRepo.UpdateUserPassword(c.app.Ctx, id, hashedPass)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (c UserController) deleteUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id, ok_id := application.GetIntVal(strings.TrimSpace(r.FormValue("userId")))

	if !ok_id {
		http.Error(rw, "Ошибка удаление пользователя", http.StatusBadRequest)
		return
	}

	user, err := c.app.Repo.UserRepo.Get(c.app.Ctx, id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Login == "admin" {
		http.Error(rw, "Ошибка удаления admin. Администратора нельзя удалить!", http.StatusBadRequest)
		return
	}

	err = c.app.Repo.AzsRepo.RemoveUserFromAzsAll(c.app.Ctx, id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.app.Repo.UserRepo.Delete(c.app.Ctx, id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (c UserController) logout(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1}
		http.SetCookie(rw, &c)
	}
	http.Redirect(rw, r, "/login", http.StatusSeeOther)
}

func (c UserController) signup(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := strings.TrimSpace(r.FormValue("name"))
	surname := strings.TrimSpace(r.FormValue("surname"))
	login := strings.TrimSpace(r.FormValue("login"))
	password := strings.TrimSpace(r.FormValue("password"))
	password2 := strings.TrimSpace(r.FormValue("password2"))
	if name == "" || surname == "" || login == "" || password == "" {
		c.signupPage(rw, "Все поля должны быть заполнены!")
		return
	}
	if password != password2 {
		c.signupPage(rw, "Пароли не совпадают! Попробуйте еще")
		return
	}
	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])
	err := c.app.Repo.UserRepo.Add(c.app.Ctx, name, surname, login, hashedPass)
	if err != nil {
		c.signupPage(rw, fmt.Sprintf("Ошибка создания пользователя: %v", err))
		return
	}
	http.Redirect(rw, r, "/users", http.StatusSeeOther)
}

func (c UserController) signupPage(rw http.ResponseWriter, message string) {

	data := Answer{message}
	sp := filepath.Join("public", "html", "signup.html")
	navi := filepath.Join("public", "html", "admin_navi.html")
	tmpl := template.Must(template.ParseFiles(sp, navi))
	err := tmpl.ExecuteTemplate(rw, "signup", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (c UserController) login(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	if login == "" || password == "" {
		c.loginPage(rw, "Необходимо указать логин и пароль!")
		return
	}
	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])
	user, err := c.app.Repo.UserRepo.Login(c.app.Ctx, login, hashedPass)
	if err != nil {
		c.loginPage(rw, "Вы ввели неверный логин или пароль!")
		return
	}

	time64 := time.Now().Unix()
	timeInt := fmt.Sprint(time64)
	token := login + password + timeInt
	hashToken := md5.Sum([]byte(token))
	hashedToken := hex.EncodeToString(hashToken[:])
	c.app.Cache[hashedToken] = user
	livingTime := 60 * time.Minute
	expiration := time.Now().Add(livingTime)

	cookie := http.Cookie{Name: "token", Value: url.QueryEscape(hashedToken), Expires: expiration}
	http.SetCookie(rw, &cookie)
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (c UserController) loginPage(rw http.ResponseWriter, message string) {
	lp := filepath.Join("public", "html", "login.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	data := Answer{message}
	err = tmpl.ExecuteTemplate(rw, "login", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
