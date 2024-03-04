package main

import (
	"api_service/heartbeat"
	"api_service/locate"
	"api_service/objects"
	"log"
	"net/http"
	"os"
)

func main() {

	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
