package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func predictHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Décodage de la requête JSON
	var req RecommendationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Requête JSON invalide", http.StatusBadRequest)
		return
	}

	resultFactor := network.Query(req.Target, req.Evidence)

	targetNode, exists := network.Nodes[req.Target]
	if !exists {
		http.Error(w, "Nœud cible introuvable dans le réseau", http.StatusNotFound)
		return
	}
	predictions := make(map[string]float64)
	for i, stateName := range targetNode.States {
		predictions[stateName] = resultFactor.Values[i]
	}

	// 4. Envoi de la réponse JSON
	resp := RecommendationResponse{
		Target:      req.Target,
		Predictions: predictions,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("ERREUR ENCODAGE JSON: %v", err)
	}
}
