package main

// RecommendationRequest est ce que le frontend envoie
type RecommendationRequest struct {
	Evidence map[string]int `json:"evidence"` 
	Target   string         `json:"target"`   
}

type RecommendationResponse struct {
	Target      string             `json:"target"`
	Predictions map[string]float64 `json:"predictions"` 
}

