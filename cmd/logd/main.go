package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	var state State
	// read statefile.

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":d", port)
	http.ListenAndServe(addr, logrot.Handler())
}
