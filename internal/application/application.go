package application

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
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

type answer struct {
	Status string `json:"status"`
	Msg    string `json:"Msg"`
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

	//TODO: защита от других пользователей
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

func (a app) ShowUsersPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	lp := filepath.Join("public", "html", "users_page.html")
	navi := filepath.Join("public", "html", "admin_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))

	users, err := a.repo.GetUserAll(a.ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = tmpl.ExecuteTemplate(rw, "User", users)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) AddUserToAsz(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id_azs, _ := getIntVal(r.FormValue("id_azs"))

	id_user, _ := getIntVal(r.FormValue("user"))
	fmt.Println(id_azs)
	fmt.Println(id_user)

	err := a.repo.AddAzsToUser(a.ctx, id_user, id_azs)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(rw, r, "/", http.StatusSeeOther)
}

func (a app) ShowHistoryReceiptsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fromSearchDate := r.FormValue("formSearch")
	toSearchDate := r.FormValue("toSearch")

	// TODO: add checking fo date from < to
	// Parse the date string
	fromSearchTime, err := time.Parse("2006-01-02", fromSearchDate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the date string
	toSearchTime, err := time.Parse("2006-01-02", toSearchDate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	a.HistoryReceiptsPage(rw, r, p, fromSearchTime, toSearchTime)
}

func (a app) HistoryReceiptsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, fromSearchTime, toSearchTime time.Time) {
	// user := r.Context().Value("user").(*repository.User)

	id_azs, ok := getIntVal(r.FormValue("id_azs"))

	if ok != true {
		http.Error(rw, "Ошибка id_azs"+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	receipts, err := a.repo.GetAzsReceiptInRange(a.ctx, id_azs, fromSearchTime.Unix(), toSearchTime.Unix())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "azs_receipt.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	azs, err := a.repo.GetAzs(a.ctx, id_azs)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	type AzsReceiptDatas struct {
		Azs           repository.AzsStatsData
		FormSearchVal string
		ToSearchVal   string
		Receipts      []repository.AzsReceiptData
		Count         int
	}

	azsReceiptDatas := AzsReceiptDatas{
		Azs:           azs,
		FormSearchVal: fromSearchTime.Format("2006-01-02"),
		ToSearchVal:   toSearchTime.Format("2006-01-02"),
		Receipts:      receipts,
		Count:         len(receipts),
	}

	err = tmpl.ExecuteTemplate(rw, "AzsReceiptDatas", azsReceiptDatas)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) AzsStats(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	idInt, ok := getIntVal(id)
	t := time.Now()
	name := strings.TrimSpace(r.FormValue("name"))
	address := strings.TrimSpace(r.FormValue("address"))
	count_colum, ok_count_colum := getIntVal(strings.TrimSpace(r.FormValue("count_colum")))
	stats := strings.TrimSpace(r.FormValue("stats"))
	fmt.Println(stats)
	rw.Header().Set("Content-Type", "application/json")
	answerStat := answer{Msg: "Ok"}
	if !ok || !ok_count_colum || id == "" || name == "" || address == "" || stats == "" {
		answerStat = answer{Msg: "error", Status: "Все поля должны быть заполнены!"}
	} else {
		azs, err := a.repo.GetAzs(a.ctx, idInt)

		if azs.Id == -1 {
			err = a.repo.AddAzs(a.ctx, idInt, 0, count_colum, t.Format(time.RFC822), name, address, stats)

			if err == nil {
				err = a.repo.CreateAzsReceipt(a.ctx, idInt)
			}

		} else if err == nil {
			azs.Time = t.Format(time.RFC822)
			azs.CountColum = count_colum
			azs.Name = name
			azs.Address = address
			azs.Stats = stats
			err = a.repo.UpdateAzsStats(a.ctx, azs)
		}

		if err != nil {
			answerStat.Status = "error"
			answerStat.Msg = err.Error()
		} else {
			answerStat = answer{"Ok", "Ok"}
		}

	}
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(answerStat)
}

func (a app) AzsReceipt(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, ok_id := getIntVal(strings.TrimSpace(r.FormValue("id")))
	time, ok_time := getIntVal(strings.TrimSpace(r.FormValue("time")))
	receipt := strings.TrimSpace(r.FormValue("receipt"))

	answerStat := answer{Msg: "Ok"}
	if ok_time != true || ok_id != true || receipt == "" {
		answerStat = answer{Msg: "error", Status: "Все поля должны быть заполнены!"}
	} else {
		err := a.repo.AddAzsReceipt(a.ctx, id, time, receipt)

		if err != nil {
			answerStat.Status = "error"
			answerStat.Msg = err.Error()
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(answerStat)
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

func (a app) Logout(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	for _, v := range r.Cookies() {
		c := http.Cookie{
			Name:   v.Name,
			MaxAge: -1}
		http.SetCookie(rw, &c)
	}
	http.Redirect(rw, r, "/login", http.StatusSeeOther)
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

func (a app) ShowUsersAzsPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {

	userId, ok := getIntVal(r.FormValue("user"))

	if !ok {
		http.Error(rw, "Error user", http.StatusBadRequest)
		return
	}
	u, err := a.repo.GetUser(a.ctx, userId)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	a.UserPage(rw, r, p, u)
}

func (a app) UserPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, u repository.User) {

	azs_statses, err := a.repo.GetAzsAllForUser(a.ctx, u.Id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "azs_stats.html")
	navi := filepath.Join("public", "html", "user_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	azses := []repository.AzsStatsDataFull{}

	for _, azs_stats := range azs_statses {

		azsStatsDataFull, err := repository.ParseStats(azs_stats)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		azses = append(azses, azsStatsDataFull)
	}

	type AzsStatsTemplate struct {
		User  repository.User
		Azses []repository.AzsStatsDataFull
	}

	azsStatsTemplate := AzsStatsTemplate{
		User:  u,
		Azses: azses,
	}

	err = tmpl.ExecuteTemplate(rw, "AzsStatsTemplate", azsStatsTemplate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) AdminPage(rw http.ResponseWriter, r *http.Request, p httprouter.Params, u repository.User, id int) {

	azs_statses, err := a.repo.GetAzsAllForUser(a.ctx, id)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "admin_page.html")
	navi := filepath.Join("public", "html", "admin_navi.html")
	tmpl := template.Must(template.ParseFiles(lp, navi))

	azses := []repository.AzsStatsDataFull{}

	for _, azs_stats := range azs_statses {

		azsStatsDataFull, err := repository.ParseStats(azs_stats)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		azses = append(azses, azsStatsDataFull)
	}

	users, err := a.repo.GetUserAll(a.ctx)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	type AdminPageTemplate struct {
		User           repository.User
		Users          []repository.User
		Azses          []repository.AzsStatsDataFull
		SelectedUserId int
	}

	adminPageTemplate := AdminPageTemplate{
		User:           u,
		Users:          users,
		Azses:          azses,
		SelectedUserId: id,
	}

	err = tmpl.ExecuteTemplate(rw, "AdminPageTemplate", adminPageTemplate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
