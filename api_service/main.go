package main

import (
	"api_service/heartbeat"
	"api_service/locate"
	"api_service/objects"
	"api_service/versions"
	"log"
	"net/http"
	"os"
)

func main() {

	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
