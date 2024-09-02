package application

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/Vadosss63/t-azs/internal/repository/user"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type App struct {
	Ctx   context.Context
	Repo  *repository.Repository
	Cache map[string]user.User
	Token string
	Port  int
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool, token string, port int) *App {
	return &App{ctx, repository.NewRepository(dbpool), make(map[string]user.User), token, port}
}

func (a App) Routes(router *httprouter.Router) {
	router.ServeFiles("/public/*filepath", http.Dir("public"))

	router.ServeFiles("/install/*filepath", http.Dir("/tmp/t_azs/update"))

	router.GET("/", a.Authorized(a.startPage))
	router.GET("/azs/control", a.Authorized(a.azsPage))

	router.GET("/azs_receipt/history", a.Authorized(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		now := time.Now()
		loc := now.Location()
		paymentType := ""

		fromSearchDateTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		toSearchDateTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

		a.historyReceiptsPage(rw, r, p, fromSearchDateTime, toSearchDateTime, paymentType)
	}))

	router.POST("/azs_receipt/history", a.Authorized(a.showHistoryReceiptsPage))

	router.POST("/azs_stats", a.Authorized(a.azsStats))
	router.DELETE("/azs_stats", a.Authorized(a.deleteAsz))

	router.POST("/azs_receipt", a.Authorized(a.azsReceipt))

	router.POST("/get_azs_button", a.Authorized(a.getAzsButton))

	router.POST("/reset_azs_button", a.Authorized(a.resetAzsButton))

	router.GET("/reset_azs_button", a.Authorized(a.resetAzs))
	router.POST("/push_azs_button", a.Authorized(a.pushAzsButton))
	router.GET("/azs_button_ready", a.Authorized(a.azsButtonReady))

	router.POST("/add_user_to_asz", a.Authorized(a.addUserToAsz))

	router.GET("/users", a.Authorized(a.showUsersPage))

	router.GET("/show_for_user", a.Authorized(a.showUsersAzsPage))

	router.POST("/show_azs_for", a.Authorized(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id_user, ok_id := GetIntVal(r.FormValue("user"))

		userId, ok := r.Context().Value("userId").(int)

		if !ok || !ok_id {
			http.Error(rw, "Error user", http.StatusBadRequest)
			return
		}
		u, err := a.Repo.UserRepo.Get(a.Ctx, userId)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		a.adminPage(rw, r, p, u, id_user)
	}))

	router.POST("/save-point", a.Authorized(a.savePointHandler))
	router.GET("/points", a.Authorized(a.pointsHandler))

}

func (a App) startPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	userId, ok := r.Context().Value("userId").(int)

	if !ok {
		http.Error(rw, "Error user", http.StatusBadRequest)
		return
	}
	u, err := a.Repo.UserRepo.Get(a.Ctx, userId)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if u.Login == "admin" {
		a.adminPage(rw, r, p, u, -1)
		return
	}

	a.userPage(rw, r, p, u)
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

func (a App) validateToken(rw http.ResponseWriter, tokenReq string) bool {
	tokenReq = strings.TrimSpace(tokenReq)
	if a.Token != tokenReq {
		return false
	}
	return true
}

func (a App) Authorized(next httprouter.Handle) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token, err := readCookie("token", r)
		if err == nil {
			if user, ok := a.Cache[token]; ok {
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
			SendJsonResponse(rw, http.StatusUnauthorized, "Invalid token", "Error")
			return
		}

		http.Redirect(rw, r, "/login", http.StatusSeeOther)
	}
}
