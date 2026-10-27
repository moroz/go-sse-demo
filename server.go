package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func serveEvents(w http.ResponseWriter, r *http.Request) {
	rc := http.NewResponseController(w)
	rc.SetWriteDeadline(time.Time{})

	w.Header().Add("Content-Type", "text/event-stream")
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("X-Accel-Buffering", "no") // disable buffering in Nginx
	if err := rc.Flush(); err != nil {
		return
	}

	clientGone := r.Context().Done()

	t := time.Tick(time.Second)

	for {
		select {
		case <-clientGone:
			log.Print("Client disconnected")
			return
		case <-t:
			_, err := fmt.Fprintf(w, "data: %s\n\n", time.Now().Format(time.RFC3339))
			if err != nil {
				return
			}
			err = rc.Flush()
			if err != nil {
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("GET /events", serveEvents)

	log.Printf("Listening on port :4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
