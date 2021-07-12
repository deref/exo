package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/deref/exo/logcol"
)

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	svc := logcol.NewService()
	go svc.Collect(ctx, &logcol.CollectInput{})
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		os.Exit(0)
	}()
	http.ListenAndServe(":"+port, logcol.NewMux("/", svc))
}
