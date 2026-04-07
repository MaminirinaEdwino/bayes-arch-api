package main

import (
	"encoding/json"
	"log"
	"net/http" // Ton package précieux
	"os"

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

func setupStructure(net *gobayes.Network) {
	if net == nil {
        log.Fatal("Le réseau passé à setupStructure est nil !")
    }
    // On définit uniquement les noms et les états possibles
    net.AddNode("TempsReel", []string{"Non", "Oui"})
    net.AddNode("Equipe", []string{"Solo", "Grande"})
    net.AddNode("Stack", []string{"PHP_Symfony", "Go_Gin", "Node_Express"})

    // On définit les liens de causalité
    net.AddEdge("TempsReel", "Stack")
    net.AddEdge("Equipe", "Stack")
    
    // Note : On ne fait PAS de SetProbabilities() ici !
    // C'est la fonction syncNetworkRules qui va le faire automatiquement.
}

func syncNetworkRules(net *gobayes.Network, rulesPath string) error {
	file, err := os.Open(rulesPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var data struct {
		StackRules []gobayes.ScoreRule `json:"stack_rules"`
	}
	
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	// On récupère le nœud qu'on veut automatiser
	stackNode := net.Nodes["Stack"]
	if stackNode != nil {
		// Magie : Le générateur calcule la table CPD complexe pour nous
		stackNode.GenerateAutomatedCPD(data.StackRules)
	}

	return nil
}

func main() {
	// 1. Charger le réseau au démarrage
	network = gobayes.NewNetwork()
	setupStructure(network)
	var err error
	// network, err = gobayes.LoadFromFile("config/architecture.json")
	err = syncNetworkRules(network, "config/rules.json")
    if err != nil {
        log.Fatal("Erreur lors de la génération des connaissances :", err)
    }

	// 2. Définir la route
	http.HandleFunc("/predict", enableCORS(predictHandler))

	// 3. Lancer le serveur
	log.Println("Serveur SEAL démarré sur :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// func predictHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var req RecommendationRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Calculer l'inférence via ton moteur gobayes
// 	resultFactor := network.Query(req.Target, req.Evidence)

// 	// Transformer le facteur en map lisible pour le JSON
// 	targetNode := network.Nodes[req.Target]
// 	predictions := make(map[string]float64)
// 	for i, stateName := range targetNode.States {
// 		predictions[stateName] = resultFactor.Values[i]
// 	}

// 	resp := RecommendationResponse{
// 		Target:      req.Target,
// 		Predictions: predictions,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(resp)
// }

func predictHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Décodage de la requête JSON
    var req RecommendationRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Requête JSON invalide", http.StatusBadRequest)
        return
    }

    // 2. Lancement de l'inférence Bayésienne
    // On demande au réseau de calculer la probabilité de 'Target'
    // sachant les 'Evidence' fournies par l'utilisateur.
    resultFactor := network.Query(req.Target, req.Evidence)

    // 3. Préparation de la réponse lisible
    // On récupère le nœud cible pour faire correspondre les noms des états (ex: "Go")
    // avec leurs probabilités respectives calculées.
    targetNode, exists := network.Nodes[req.Target]
    if !exists {
        http.Error(w, "Nœud cible introuvable dans le réseau", http.StatusNotFound)
        return
    }

    predictions := make(map[string]float64)
    for i, stateName := range targetNode.States {
        // resultFactor.Values contient les probabilités normalisées
        predictions[stateName] = resultFactor.Values[i]
    }

    // 4. Envoi de la réponse JSON
    resp := RecommendationResponse{
        Target:      req.Target,
        Predictions: predictions,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}