package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type MonitorHandler struct {
	ctx    context.Context
	events chan []byte
}

func (h MonitorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("INFO: %s %s %s\n", r.Method, r.URL.Path, r.RemoteAddr)

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	flusher, ok := w.(http.Flusher)

	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case msg, ok := <-h.events:
			if !ok {
				http.Error(w, "streaming stopped", http.StatusInternalServerError)
				return
			}

			id := fmt.Sprint(time.Now().Unix())

			_, _ = w.Write(NewEvent(id, "message", string(msg)).Buffer().Bytes())
			flusher.Flush()

		case <-h.ctx.Done():
			return

		case <-r.Context().Done():
			log.Printf("INFO: client %s disconnected\n", r.RemoteAddr)
			return
		}
	}
}
