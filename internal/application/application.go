package application

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/Vadosss63/t-azs/internal/repository/user"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type Settings struct {
	Port     int    `json:"port"`
	Token    string `json:"token"`
	YaApiKey string `json:"ya_api_key"`
	YaPayUrl string `json:"ya_pay_url"`
}

type App struct {
	Ctx      context.Context
	Repo     *repository.Repository
	Cache    map[string]user.User
	Settings Settings
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool, settings Settings) *App {
	return &App{ctx, repository.NewRepository(dbpool), make(map[string]user.User), settings}
}

func (a App) Routes(router *httprouter.Router) {
	router.ServeFiles("/public/*filepath", http.Dir("public"))

	router.ServeFiles("/install/*filepath", http.Dir("/tmp/t_azs/update"))

}

func (a App) GetYaPayApiKey() string {
	return a.Settings.YaApiKey
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

func (a App) ValidateToken(tokenReq string) bool {
	tokenReq = strings.TrimSpace(tokenReq)
	return a.Settings.Token == tokenReq
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
			if a.ValidateToken(tokenReq) {
				next(rw, r, ps)
				return
			}
			SendJsonResponse(rw, http.StatusUnauthorized, "Invalid token", "Error")
			return
		}

		http.Redirect(rw, r, "/login", http.StatusSeeOther)
	}
}
