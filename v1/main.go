package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"oss_start/v1/objects"
)

func main() {

	// 创建目录
	//err := os.MkdirAll("./tmp/objects", 0755)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//_, err = os.Create("./tmp/objects/test")
	//if err != nil {
	//	log.Println(err)
	//}
	http.HandleFunc("/objects/", objects.Handler)
	fmt.Println(os.Getenv("STORAGE_ROOT"))
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
