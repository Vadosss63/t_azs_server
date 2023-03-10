package application

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (a app) ResetPasswordUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id, ok_id := getIntVal(strings.TrimSpace(r.FormValue("userId")))
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

	err := a.repo.UpdateUserPassword(a.ctx, id, hashedPass)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func (a app) DeleteUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id, ok_id := getIntVal(strings.TrimSpace(r.FormValue("userId")))

	if !ok_id {
		http.Error(rw, "Ошибка удаление пользователя", http.StatusBadRequest)
		return
	}

	err := a.repo.DeleteUser(a.ctx, id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) Logout(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1}
		http.SetCookie(rw, &c)
	}
	http.Redirect(rw, r, "/login", http.StatusSeeOther)
}

func readCookie(name string, r *http.Request) (value string, err error) {
	if name == "" {
		return value, errors.New("you are trying to read empty cookie")
	}
	cookie, err := r.Cookie(name)
	if err != nil {
		return value, err
	}
	str := cookie.Value
	value, _ = url.QueryUnescape(str)
	return value, err
}

func (a app) Signup(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := strings.TrimSpace(r.FormValue("name"))
	surname := strings.TrimSpace(r.FormValue("surname"))
	login := strings.TrimSpace(r.FormValue("login"))
	password := strings.TrimSpace(r.FormValue("password"))
	password2 := strings.TrimSpace(r.FormValue("password2"))
	if name == "" || surname == "" || login == "" || password == "" {
		a.SignupPage(rw, "Все поля должны быть заполнены!")
		return
	}
	if password != password2 {
		a.SignupPage(rw, "Пароли не совпадают! Попробуйте еще")
		return
	}
	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])
	err := a.repo.AddNewUser(a.ctx, name, surname, login, hashedPass)
	if err != nil {
		a.SignupPage(rw, fmt.Sprintf("Ошибка создания пользователя: %v", err))
		return
	}
	http.Redirect(rw, r, "/users", http.StatusSeeOther)
}

func (a app) SignupPage(rw http.ResponseWriter, message string) {
	sp := filepath.Join("public", "html", "signup.html")
	navi := filepath.Join("public", "html", "admin_navi.html")
	tmpl := template.Must(template.ParseFiles(sp, navi))

	type answer struct {
		Message string
	}
	data := answer{message}
	err := tmpl.ExecuteTemplate(rw, "signup", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) Login(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	login := r.FormValue("login")
	password := r.FormValue("password")
	if login == "" || password == "" {
		a.LoginPage(rw, "Необходимо указать логин и пароль!")
		return
	}
	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])
	user, err := a.repo.Login(a.ctx, login, hashedPass)
	if err != nil {
		a.LoginPage(rw, "Вы ввели неверный логин или пароль!")
		return
	}
	//логин и пароль совпадают, поэтому генерируем токен, пишем его в кеш и в куки
	time64 := time.Now().Unix()
	timeInt := string(time64)
	token := login + password + timeInt
	hashToken := md5.Sum([]byte(token))
	hashedToken := hex.EncodeToString(hashToken[:])
	a.cache[hashedToken] = user
	livingTime := 60 * time.Minute
	expiration := time.Now().Add(livingTime)
	//кука будет жить 1 час
	cookie := http.Cookie{Name: "token", Value: url.QueryEscape(hashedToken), Expires: expiration}
	http.SetCookie(rw, &cookie)
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (a app) LoginPage(rw http.ResponseWriter, message string) {
	lp := filepath.Join("public", "html", "login.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	type answer struct {
		Message string
	}
	data := answer{message}
	err = tmpl.ExecuteTemplate(rw, "login", data)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) Authorized(next httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := readCookie("token", r)
		if err != nil {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}

		user, ok := a.cache[token]
		if !ok {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}
		// Call the next handler with the user information
		next(rw, r.WithContext(context.WithValue(r.Context(), "userId", user.Id)), ps)
	}
}
