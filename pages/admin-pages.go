package pages

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"
	"web/models"
)

func AdminMain(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/admin-main.html",
		"templates/admin-header.html",
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

	tmpl.ExecuteTemplate(w, "admin-main", data)
}

func PanicPage(w http.ResponseWriter, r *http.Request) {
	panic("this must me recovered")
}

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("adminAuthMiddleware", r.URL.Path)
		_, err := r.Cookie("session_id")
		fmt.Println("check Auth")
		if err != nil {
			fmt.Println("no auth at", r.URL.Path)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	if err == http.ErrNoCookie {
		fmt.Println("no Cookie")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	http.Redirect(w, r, "/", http.StatusFound)
}
