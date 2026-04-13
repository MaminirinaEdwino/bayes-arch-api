package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/MaminirinaEdwino/gobayes"
)

func setupStructure(net *gobayes.Network) {
	if net == nil {
		log.Fatal("Le réseau passé à setupStructure est nil !")
	}
	net.AddNode("TempsReel", []string{"Non", "Oui"})
    net.AddNode("Equipe", []string{"Solo", "Grande"})
    net.AddNode("Complexite", []string{"Faible", "Elevee"})


    net.AddNode("Stack", []string{"PHP_Symfony", "Go_Gin", "Node"})

    net.AddEdge("TempsReel", "Complexite")
    net.AddEdge("Equipe", "Complexite")
    net.AddEdge("Complexite", "Stack")
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
