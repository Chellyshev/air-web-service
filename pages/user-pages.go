package pages

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
	"web/models"
)

type MainPageData struct {
	AQI    int
	Danger string
	Level  string
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/client-main.html",
		"templates/header.html",
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	api_data := new(models.JsonPrediction)
	resp, err := http.Get("http://127.0.0.1:8080/forecast")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	bytes_data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.Unmarshal(bytes_data, api_data)
	stringsData := strings.Split(api_data.LLMResponse, "\n")
	if len(stringsData) <= 1 {
		http.Error(w, "I haven't data", 500)
	}
	data := MainPageData{AQI: 60,
		Danger: stringsData[1],
		Level:  stringsData[10]}

	tmpl.ExecuteTemplate(w, "client-main", data)

}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/login.html",
		"templates/header.html",
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	tmpl.ExecuteTemplate(w, "login", nil)
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	r.ParseForm()
	hash := md5.New()
	login := r.FormValue("login")
	io.WriteString(hash, login)

	fmt.Println(r.FormValue("password"))
	expiration := time.Now().Add(10 * time.Hour)
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   fmt.Sprintf("%x", hash.Sum(nil)),
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/admin/", http.StatusSeeOther)
}
