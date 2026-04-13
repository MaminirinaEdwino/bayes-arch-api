package main

import (
	"log"

	"github.com/MaminirinaEdwino/gobayes"
	"github.com/fsnotify/fsnotify"
)

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
                if event.Has(fsnotify.Write) {
                    log.Println("🔄 Modification détectée dans ",filename,", rechargement...")
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
    
    <-make(chan struct{})
}