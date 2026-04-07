package main

import (
	"encoding/json"
	"log"
	"net/http" // Ton package précieux

	"github.com/MaminirinaEdwino/gobayes"
)

var network *gobayes.Network

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Si c'est une requête de pré-vérification (OPTIONS), on s'arrête là
		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	}
}

func main() {
	// 1. Charger le réseau au démarrage
	var err error
	network, err = gobayes.LoadFromFile("config/architecture.json")
	if err != nil {
		log.Fatal("Impossible de charger le réseau :", err)
	}

	// 2. Définir la route
	http.HandleFunc("/predict", enableCORS(predictHandler))

	// 3. Lancer le serveur
	log.Println("Serveur SEAL démarré sur :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func predictHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var req RecommendationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculer l'inférence via ton moteur gobayes
	resultFactor := network.Query(req.Target, req.Evidence)

	// Transformer le facteur en map lisible pour le JSON
	targetNode := network.Nodes[req.Target]
	predictions := make(map[string]float64)
	for i, stateName := range targetNode.States {
		predictions[stateName] = resultFactor.Values[i]
	}

	resp := RecommendationResponse{
		Target:      req.Target,
		Predictions: predictions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}