package models

type JsonPrediction struct {
	Prediction  map[string]map[int]float64 `json:"prediction"`
	Decision    map[string]any             `json:"decision"`
	LLMResponse string                     `json:"llm_response"`
}

type LLMPrediction struct {
	Situation, Warning  string
	Pm25, Pm10, CO, NO2 float64
	ToDo, Facts         []string
}

type Breakpoint struct {
	Clo float64
	Chi float64
	Ilo int
	Ihi int
}

var Pm25Table = []Breakpoint{
	{0, 12, 0, 50},
	{12.1, 35.4, 51, 100},
	{35.5, 55.4, 101, 150},
	{55.5, 150, 151, 200},
	{150.1, 250, 201, 300},
	{250.1, 500, 301, 500},
}

func CalcAQI(c float64, table []Breakpoint) int {
	for _, b := range table {
		if c >= b.Clo && c <= b.Chi {
			aqi := float64(b.Ihi-b.Ilo)/(b.Chi-b.Clo)*(c-b.Clo) + float64(b.Ilo)
			return int(aqi)
		}
	}
	return 0
}
