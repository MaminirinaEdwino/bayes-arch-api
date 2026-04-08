package main

import (
	"log"

	"github.com/MaminirinaEdwino/gobayes"
	"github.com/fsnotify/fsnotify"
)

// WatchRules surveille les changements sur le fichier rules.json
func watchRules(filename string, net *gobayes.Network) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()

    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                // On réagit uniquement si le fichier est écrit (modifié)
                if event.Has(fsnotify.Write) {
                    log.Println("🔄 Modification détectée dans ",filename,", rechargement...")
                    
                    // On recharge les règles et on régénère la CPD
                    syncNetworkRules(net, filename) 
                    
                    log.Println("✅ Réseau mis à jour avec succès !")
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("Erreur watcher:", err)
            }
        }
    }()

    err = watcher.Add(filename)
    if err != nil {
        log.Fatal(err)
    }
    
    // Bloquer la goroutine pour qu'elle continue de tourner
    <-make(chan struct{})
}