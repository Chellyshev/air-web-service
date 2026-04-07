package pages

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"web/models"
)

type MainPageData struct {
	AQI       int
	Warning   string
	Situation string
}

var Db *sql.DB
var APIData = new(models.JsonPrediction)
var LLMPrediction = new(models.LLMPrediction)

func Scheduler() {
	go func() {
		Update()
		for {
			time.Sleep(24 * time.Hour)
			Update()
		}
	}()
}

func cleanBullet(line string) string {
	line = strings.TrimPrefix(line, "-")
	line = strings.TrimSpace(line)
	return line
}

func parseValue(line string) {
	parts := strings.Split(line, "—")
	if len(parts) < 2 {
		return
	}

	name := strings.TrimSpace(parts[0])
	valuePart := strings.TrimSpace(parts[1])

	valueStr := strings.Split(valuePart, "(")[0]
	valueStr = strings.TrimSpace(valueStr)

	val, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return
	}

	switch {
	case strings.Contains(name, "PM2.5"):
		LLMPrediction.Pm25 = math.Round(val*100) / 100
	case strings.Contains(name, "PM10"):
		LLMPrediction.Pm10 = math.Round(val*100) / 100
	case strings.Contains(name, "NO2"):
		LLMPrediction.NO2 = math.Round(val*100) / 100
	case strings.Contains(name, "CO"):
		LLMPrediction.CO = math.Round(val*100) / 100
	}
}

func Update() {
	resp, err := http.Get("http://127.0.0.1:8080/forecast")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	bytes_data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	json.Unmarshal(bytes_data, APIData)
	lines := strings.Split(APIData.LLMResponse, "\n")

	section := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "СИТУАЦИЯ"):
			section = "situation"
			continue
		case strings.HasPrefix(line, "ПОКАЗАТЕЛИ"):
			section = "values"
			continue
		case strings.HasPrefix(line, "РИСК"):
			section = "risk"
			continue
		case strings.HasPrefix(line, "ФАКТОРЫ"):
			section = "facts"
			continue
		case strings.HasPrefix(line, "МЕРЫ"):
			section = "todo"
			continue
		}

		switch section {
		case "situation":
			LLMPrediction.Situation += line + " "
		case "risk":
			LLMPrediction.Warning += line + " "
		case "values":
			parseValue(line)
		case "facts":
			LLMPrediction.Facts = append(LLMPrediction.Facts, cleanBullet(line))
		case "todo":
			LLMPrediction.ToDo = append(LLMPrediction.ToDo, cleanBullet(line))
		}
	}

	fmt.Println(LLMPrediction.ToDo)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_id")
	if err != http.ErrNoCookie {
		http.Redirect(w, r, "/admin/", http.StatusFound)
		return
	}
	tmpl, err := template.ParseFiles(
		"templates/client-main.html",
		"templates/header.html",
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := MainPageData{AQI: models.CalcAQI(LLMPrediction.Pm25, models.Pm25Table),
		Warning:   LLMPrediction.Warning,
		Situation: LLMPrediction.Situation}
	allData := struct {
		Data   MainPageData
		Active string
	}{Data: data, Active: "main"}
	tmpl.ExecuteTemplate(w, "client-main", allData)

}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/login.html",
		"templates/header.html",
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	data := struct{ Active string }{"login"}
	tmpl.ExecuteTemplate(w, "login", data)
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

func Charts(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/charts.html",
		"templates/header.html",
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	data := struct{ Active string }{"charts"}
	tmpl.ExecuteTemplate(w, "charts", data)
}
func GetWeekData(w http.ResponseWriter, r *http.Request) {
	rows, err := Db.Query(`
		SELECT time, pm2p5, pm10, no2, co
		FROM monthly_air_data
		WHERE time >= NOW() - INTERVAL '7 days'
		ORDER BY time
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	type Point struct {
		Time string  `json:"time"`
		PM25 float64 `json:"pm25"`
		PM10 float64 `json:"pm10"`
		NO2  float64 `json:"no2"`
		CO   float64 `json:"co"`
	}

	var data []Point

	for rows.Next() {
		var t time.Time
		var p Point

		rows.Scan(&t, &p.PM25, &p.PM10, &p.NO2, &p.CO)
		p.Time = t.Format("2006-01-02 15:04")

		data = append(data, p)
	}

	json.NewEncoder(w).Encode(data)
}
