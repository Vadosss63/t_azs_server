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

type Answer struct {
	Message string
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

func (a app) resetPasswordUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

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

func (a app) deleteUser(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id, ok_id := getIntVal(strings.TrimSpace(r.FormValue("userId")))

	if !ok_id {
		http.Error(rw, "Ошибка удаление пользователя", http.StatusBadRequest)
		return
	}

	user, err := a.repo.GetUser(a.ctx, id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Login == "admin" {
		http.Error(rw, "Ошибка удаления admin. Администратора нельзя удалить!", http.StatusBadRequest)
		return
	}

	err = a.repo.RemoveUserFromAzsAll(a.ctx, id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.repo.DeleteUser(a.ctx, id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) logout(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1}
		http.SetCookie(rw, &c)
	}
	http.Redirect(rw, r, "/login", http.StatusSeeOther)
}

func (a app) signup(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	name := strings.TrimSpace(r.FormValue("name"))
	surname := strings.TrimSpace(r.FormValue("surname"))
	login := strings.TrimSpace(r.FormValue("login"))
	password := strings.TrimSpace(r.FormValue("password"))
	password2 := strings.TrimSpace(r.FormValue("password2"))
	if name == "" || surname == "" || login == "" || password == "" {
		a.signupPage(rw, "Все поля должны быть заполнены!")
		return
	}
	if password != password2 {
		a.signupPage(rw, "Пароли не совпадают! Попробуйте еще")
		return
	}
	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])
	err := a.repo.AddNewUser(a.ctx, name, surname, login, hashedPass)
	if err != nil {
		a.signupPage(rw, fmt.Sprintf("Ошибка создания пользователя: %v", err))
		return
	}
	http.Redirect(rw, r, "/users", http.StatusSeeOther)
}

func (a app) signupPage(rw http.ResponseWriter, message string) {

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

func (a app) login(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	if login == "" || password == "" {
		a.loginPage(rw, "Необходимо указать логин и пароль!")
		return
	}
	hash := md5.Sum([]byte(password))
	hashedPass := hex.EncodeToString(hash[:])
	user, err := a.repo.Login(a.ctx, login, hashedPass)
	if err != nil {
		a.loginPage(rw, "Вы ввели неверный логин или пароль!")
		return
	}

	time64 := time.Now().Unix()
	timeInt := string(time64)
	token := login + password + timeInt
	hashToken := md5.Sum([]byte(token))
	hashedToken := hex.EncodeToString(hashToken[:])
	a.cache[hashedToken] = user
	livingTime := 60 * time.Minute
	expiration := time.Now().Add(livingTime)

	cookie := http.Cookie{Name: "token", Value: url.QueryEscape(hashedToken), Expires: expiration}
	http.SetCookie(rw, &cookie)
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (a app) loginPage(rw http.ResponseWriter, message string) {
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

func (a app) validateToken(rw http.ResponseWriter, tokenReq string) bool {
	tokenReq = strings.TrimSpace(tokenReq)
	if a.token != tokenReq {
		return false
	}
	return true
}

func (a app) authorized(next httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := readCookie("token", r)
		if err == nil {
			if user, ok := a.cache[token]; ok {
				next(rw, r.WithContext(context.WithValue(r.Context(), "userId", user.Id)), ps)
				return
			}
		}

		tokenReq := r.FormValue("token")
		if tokenReq != "" {
			if a.validateToken(rw, tokenReq) {
				next(rw, r, ps)
				return
			}
			sendJsonResponse(rw, http.StatusUnauthorized, "Invalid token", "Error")
			return
		}

		http.Redirect(rw, r, "/login", http.StatusSeeOther)
	}
}
