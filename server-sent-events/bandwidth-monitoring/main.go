package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	g, gCtx := errgroup.WithContext(ctx)

	events := make(chan []byte, 1)

	mux := http.NewServeMux()
	mux.Handle("/events", MonitorHandler{ctx: gCtx, events: events})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	service := Service{
		changeTicker: time.NewTicker(1 * time.Second),
		buffer:       make(chan *NetStat, 1),
		events:       events,
	}

	g.Go(func() error {
		return service.Capture(gCtx)
	})

	g.Go(func() error {
		return service.Transform(gCtx)
	})

	g.Go(func() error {
		log.Println("INFO: server listening at", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		log.Printf("INFO: received termination signal\n")

		shutdownCtx, release := context.WithTimeout(context.Background(), 5*time.Second)
		defer release()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return err
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		log.Panicln("ERROR:", err)
	}

	log.Println("INFO:", "program exits with code 0")
}
