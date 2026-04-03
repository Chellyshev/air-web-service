package models

type JsonPrediction struct {
	Prediction  map[string]map[int]float64 `json:"prediction"`
	Decision    map[string]any             `json:"decision"`
	LLMResponse string                     `json:"llm_response"`
}
