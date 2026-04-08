package main

import (
	"encoding/json"
	"log"
	"net/http" // Ton package précieux
	"os"

	"github.com/MaminirinaEdwino/gobayes"
)

var network *gobayes.Network

func setupStructure(net *gobayes.Network) {
	if net == nil {
		log.Fatal("Le réseau passé à setupStructure est nil !")
	}
	// On définit uniquement les noms et les états possibles
	net.AddNode("TempsReel", []string{"Non", "Oui"})
	net.AddNode("Equipe", []string{"Solo", "Grande"})
	net.AddNode("Stack", []string{"PHP_Symfony", "Go_Gin", "Node"})

	// On définit les liens de causalité
	net.AddEdge("TempsReel", "Stack")
	net.AddEdge("Equipe", "Stack")
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
		stackNode.GenerateCPD(data.StackRules)
		log.Printf("CPD de la Stack générée : %d valeurs", len(stackNode.CPD))
		if len(stackNode.CPD) == 0 {
			log.Fatal("ERREUR : La CPD est vide. Vérifie tes règles dans rules.json")
		}
	}

	return nil
}

func main() {
	network = gobayes.NewNetwork()
	setupStructure(network)
	var err error
	err = syncNetworkRules(network, "config/rules.json")
	if err != nil {
		log.Fatal("Erreur lors de la génération des connaissances :", err)
	}

	http.HandleFunc("/predict", enableCORS(predictHandler))

	// 3. Lancer le serveur
	log.Println("Serveur SEAL démarré sur :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


