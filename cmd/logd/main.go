package main

import (
	"net/http"
	"os"

	"github.com/deref/exo/logrot"
)

func main() {
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, logrot.NewHandler())
}
