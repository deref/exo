package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/deref/exo/logrot"
)

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	svc := logrot.NewService()
	go svc.CollectLogs(ctx, &logrot.CollectLogsInput{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		os.Exit(0)
	}()
	http.ListenAndServe(":"+port, logrot.NewMux("/", svc))
}
