package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startServer(listen string, mux *http.ServeMux, cancel context.CancelFunc) {
	srv := &http.Server{
		Addr:    listen,
		Handler: mux,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Listening on port %s ...\n", listen)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	<-sigCh
	cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelTimeout()

	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
}
