package main

// RecommendationRequest est ce que le frontend envoie
type RecommendationRequest struct {
	Evidence map[string]int `json:"evidence"` // ex: {"TempsReel": 1, "Budget": 0}
	Target   string         `json:"target"`   // ex: "Stack"
}

// RecommendationResponse est ce que l'API renvoie
type RecommendationResponse struct {
	Target      string             `json:"target"`
	Predictions map[string]float64 `json:"predictions"` // ex: {"Go": 0.90, "PHP": 0.10}
}