package application

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Vadosss63/t-azs/internal/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
)

type app struct {
	ctx   context.Context
	repo  *repository.Repository
	cache map[string]repository.User
	token string
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool, token string) *app {
	return &app{ctx, repository.NewRepository(dbpool), make(map[string]repository.User), token}
}

func (a app) Routes(router *httprouter.Router) {
	router.ServeFiles("/public/*filepath", http.Dir("public"))

	router.GET("/", a.authorized(a.startPage))

	router.GET("/login", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		a.loginPage(rw, "")
	})

	router.POST("/login", a.login)

	router.GET("/azs/control", a.authorized(a.azsPage))

	router.GET("/logout", a.logout)

	router.GET("/signup", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		a.signupPage(rw, "")
	})

	router.GET("/azs_receipt/history", a.authorized(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		now := time.Now()
		loc := now.Location()
		paymentType := ""

		fromSearchDateTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		toSearchDateTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)

		a.historyReceiptsPage(rw, r, p, fromSearchDateTime, toSearchDateTime, paymentType)
	}))

	router.POST("/azs_receipt/history", a.authorized(a.showHistoryReceiptsPage))

	router.POST("/signup", a.signup)

	router.POST("/azs_stats", a.azsStats)
	router.DELETE("/azs_stats", a.authorized(a.deleteAsz))

	router.POST("/azs_receipt", a.azsReceipt)

	router.POST("/get_azs_button", a.getAzsButton)

	router.POST("/reset_azs_button", a.resetAzsButton)

	router.GET("/reset_azs_button", a.authorized(a.resetAzs))

	router.POST("/push_azs_button", a.authorized(a.pushAzsButton))

	router.GET("/azs_button_ready", a.azsButtonReady)

	router.POST("/get_log_cmd", a.getLogButton)
	router.POST("/upload_log", a.uploadLogs)
	router.POST("/reset_log_cmd", a.resetLogButton)
	router.GET("/set_log_cmd", a.setLogCmd)

	router.GET("/list_logs", a.listLogFiles)
	router.GET("/download_log", a.downloadLogFile)
	router.GET("/delete_logs", a.deleteLogs)

	router.POST("/get_app_update_button", a.getAppUpdateButton)
	router.POST("/reset_app_update_button", a.resetAppUpdateButton)
	router.GET("/set_app_update_cmd", a.setAppUpdateCmd)

	router.POST("/add_user_to_asz", a.authorized(a.addUserToAsz))

	router.GET("/users", a.authorized(a.showUsersPage))

	router.DELETE("/user", a.authorized(a.deleteUser))

	router.POST("/reset_password", a.authorized(a.resetPasswordUser))

	router.GET("/show_for_user", a.authorized(a.showUsersAzsPage))

	router.POST("/show_azs_for", a.authorized(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		id_user, ok_id := getIntVal(r.FormValue("user"))

		userId, ok := r.Context().Value("userId").(int)

		if !ok || !ok_id {
			http.Error(rw, "Error user", http.StatusBadRequest)
			return
		}
		u, err := a.repo.GetUser(a.ctx, userId)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		a.adminPage(rw, r, p, u, id_user)
	}))
}

func getIntVal(val string) (int, bool) {
	sum, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	return sum, true
}

func (a app) startPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	userId, ok := r.Context().Value("userId").(int)

	if !ok {
		http.Error(rw, "Error user", http.StatusBadRequest)
		return
	}
	u, err := a.repo.GetUser(a.ctx, userId)

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
