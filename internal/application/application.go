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
}

func getIntVal(val string) (int, bool) {
	sum, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	return sum, true
}

func NewApp(ctx context.Context, dbpool *pgxpool.Pool) *app {
	return &app{ctx, repository.NewRepository(dbpool), make(map[string]repository.User)}
}

func (a app) Routes(r *httprouter.Router) {
	r.ServeFiles("/public/*filepath", http.Dir("public"))

	r.GET("/", a.Authorized(a.StartPage))

	r.GET("/login", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		a.LoginPage(rw, "")
	})

	r.POST("/login", a.Login)

	r.GET("/logout", a.Logout)

	r.GET("/signup", func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		a.SignupPage(rw, "")
	})

	r.GET("/azs_receipt/history", a.Authorized(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		now := time.Now()
		oneMonthAgo := now.AddDate(0, -1, 0)
		a.HistoryReceiptsPage(rw, r, p, oneMonthAgo, now)
	}))

	r.POST("/azs_receipt/history", a.Authorized(a.ShowHistoryReceiptsPage))

	r.POST("/signup", a.Signup)

	r.POST("/azs_stats", a.AzsStats)

	r.POST("/azs_receipt", a.AzsReceipt)

	r.POST("/add_user_to_asz", a.Authorized(a.AddUserToAsz))

	r.GET("/users", a.Authorized(a.ShowUsersPage))

	r.DELETE("/user", a.Authorized(a.DeleteUser))
	r.POST("/reset_password", a.Authorized(a.ResetPasswordUser))

	r.GET("/show_for_user", a.Authorized(a.ShowUsersAzsPage))

	r.POST("/show_azs_for", a.Authorized(func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
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
		a.AdminPage(rw, r, p, u, id_user)
	}))
}

func (a app) StartPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

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
		a.AdminPage(rw, r, p, u, -1)
		return
	}

	a.UserPage(rw, r, p, u)
}
