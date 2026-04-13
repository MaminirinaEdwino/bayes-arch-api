package main

import (
	"log"
	"net/http"

	"github.com/MaminirinaEdwino/gobayes"
)

var network *gobayes.Network

func main() {
	network = gobayes.NewNetwork()
	setupStructure(network)
	var err error
	err = syncNetworkRules(network, "config/rules.json")
	if err != nil {
		log.Fatal("Erreur lors de la génération des connaissances :", err)
	}

	http.HandleFunc("/predict", enableCORS(predictHandler))
	go watchRules("config/rules.json", network)
	go watchRules("main.go", network)
	
	log.Println("Serveur SEAL démarré sur :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}


