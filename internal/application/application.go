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

	r.POST("/azs_receipt/history", a.ShowHistoryReceiptsPage)

	r.POST("/signup", a.Signup)

	r.POST("/azs_stats", a.AzsStats)

	r.POST("/azs_receipt", a.AzsReceipt)
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

	id_azs, ok := getIntVal(r.FormValue("id_azs"))

	if ok != true {
		http.Error(rw, "Ошибка id_azs"+r.FormValue("id_azs"), http.StatusBadRequest)
		return
	}

	receipts, err := a.repo.GetAzsReceiptInRange(a.ctx, id_azs-1, fromSearchTime.Unix(), toSearchTime.Unix())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "azs_receipt.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	type AzsReceiptDatas struct {
		FormSearchVal string
		ToSearchVal   string
		Receipts      []repository.AzsReceiptData
	}

	azsReceiptDatas := AzsReceiptDatas{
		FormSearchVal: fromSearchTime.Format("2006-01-02"),
		ToSearchVal:   toSearchTime.Format("2006-01-02"),
		Receipts:      receipts,
	}

	err = tmpl.ExecuteTemplate(rw, "AzsReceiptDatas", azsReceiptDatas)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}

func (a app) AzsStats(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := strings.TrimSpace(r.FormValue("id"))
	t := time.Now()
	name := strings.TrimSpace(r.FormValue("name"))
	address := strings.TrimSpace(r.FormValue("address"))
	stats := strings.TrimSpace(r.FormValue("stats"))
	fmt.Println(stats)
	rw.Header().Set("Content-Type", "application/json")
	answerStat := answer{Msg: "Ok"}
	if id == "" || name == "" || address == "" || stats == "" {
		answerStat = answer{Msg: "error", Status: "Все поля должны быть заполнены!"}
	} else {
		answerStat.Msg = id + name + address + stats

		idInt, ok := getIntVal(id)
		if ok == true {
			azs, err := a.repo.GetAzs(a.ctx, idInt)
			if err != nil {
				answerStat.Status = "error"
				answerStat.Msg = err.Error()
			}
			if azs.Id == -1 {
				err := a.repo.AddAzs(a.ctx, idInt, 0, t.Format(time.RFC822), name, address, stats)
				if err != nil {
					answerStat.Status = "error"
					answerStat.Msg = err.Error()
				}
			}
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
		//Как прокинуть функцию Сделать свой хендлер?
		if _, ok := a.cache[token]; !ok {
			http.Redirect(rw, r, "/login", http.StatusSeeOther)
			return
		}
		next(rw, r, ps)
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
	a.LoginPage(rw, fmt.Sprintf("%s, вы успешно зарегистрированы! Теперь вам доступен вход через страницу авторизации", name))
}

func (a app) SignupPage(rw http.ResponseWriter, message string) {
	sp := filepath.Join("public", "html", "signup.html")
	tmpl, err := template.ParseFiles(sp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	type answer struct {
		Message string
	}
	data := answer{message}
	err = tmpl.ExecuteTemplate(rw, "signup", data)
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
	azs_stats, err := a.repo.GetAzs(a.ctx, 10111992)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	lp := filepath.Join("public", "html", "azs_stats.html")
	tmpl, err := template.ParseFiles(lp)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	azsStatsDataFull, err := repository.ParseStats(azs_stats)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	type AzsStatsTemplate struct {
		Id    string
		Azses []repository.AzsStatsDataFull
	}

	azsStatsTemplate := AzsStatsTemplate{
		Id:    "10111992",
		Azses: []repository.AzsStatsDataFull{azsStatsDataFull, azsStatsDataFull},
	}

	err = tmpl.ExecuteTemplate(rw, "AzsStatsTemplate", azsStatsTemplate)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
}
