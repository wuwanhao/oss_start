package main

import (
	"data_service/heartbeat"
	"data_service/locate"
	"data_service/objects"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	fmt.Println(os.Getenv("STORAGE_ROOT"))
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
