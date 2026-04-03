package pages

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"strings"
	"web/models"
)

type MainPageData struct {
	AQI    int
	Danger string
	Level  string
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/client-main.html")
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

	data := MainPageData{AQI: 60,
		Danger: stringsData[1],
		Level:  stringsData[10]}

	tmpl.Execute(w, data)

}
