package pages

import (
	"database/sql"
	"fmt"
	"html/template"
	"math"
	"net/http"
	"time"
	"web/models"
)

type AirData struct {
	PM25 float64
	PM10 float64
	NO2  float64
	CO   float64
	AQI  int
}

type adminPageData struct {
	ActualAirData, PredictionAirData AirData
	Situation, Warning               string
	Facts, ToDo                      []string
}

func GetLatest(db *sql.DB) (AirData, error) {
	var d AirData

	err := db.QueryRow(`
		SELECT pm2p5, pm10, no2, co
		FROM monthly_air_data
		ORDER BY time DESC
		LIMIT 1
	`).Scan(&d.PM25, &d.PM10, &d.NO2, &d.CO)
	d.PM25 = math.Round(d.PM25*100) / 100
	d.PM10 = math.Round(d.PM10*100) / 100
	d.NO2 = math.Round(d.NO2*100) / 100
	d.CO = math.Round(d.CO*100) / 100
	return d, err
}

func AdminMain(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/admin-main.html",
		"templates/admin-header.html",
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	now, err := GetLatest(Db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	future := AirData{PM25: LLMPrediction.Pm25, PM10: LLMPrediction.Pm10, CO: LLMPrediction.CO, NO2: LLMPrediction.NO2}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	now.AQI = models.CalcAQI(now.PM25, models.Pm25Table)
	future.AQI = models.CalcAQI(LLMPrediction.Pm25, models.Pm25Table)

	data := adminPageData{
		ActualAirData:     now,
		PredictionAirData: future,
		Situation:         LLMPrediction.Situation,
		Warning:           LLMPrediction.Warning,
		Facts:             LLMPrediction.Facts,
		ToDo:              LLMPrediction.ToDo,
	}

	tmpl.ExecuteTemplate(w, "admin-main", data)
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
