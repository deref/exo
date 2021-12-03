package main

import (
	"fmt"
	"net/http"
	"os"
)

func handle(w http.ResponseWriter, req *http.Request) {
	response, _ := os.LookupEnv("RESPONSE")
	if _, err := w.Write([]byte(response)); err != nil {
		panic(err)
	}
}

func main() {
	host, hostSet := os.LookupEnv("HOST")
	if !hostSet {
		host = "0.0.0.0"
	}

	port, _ := os.LookupEnv("PORT")
	listenAddress := host + ":" + port
	fmt.Printf("listening on: %+v\n", listenAddress)

	http.HandleFunc("/", handle)
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		panic(err)
	}
}
