package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	log.Print("starting")

	urls := make(chan string)
	cancels := make([]context.CancelFunc, 0)

	playSong := func(w http.ResponseWriter, _ *http.Request) {
		for {
			log.Print("hit")
			url := <-urls
			log.Print("playing ", url)
			res, err := http.Get(url)
			if err != nil {
				log.Print(err)
				continue
			}
			ctx, cancel := context.WithCancel(context.Background())
			cancels = append(cancels, cancel)

			Copy(ctx, w, res.Body)
			log.Print("pop")

		}
	}
	addSong := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			for _, cancel := range cancels {
				cancel()
			}
			cancels = make([]context.CancelFunc, 0)
			if err := r.ParseForm(); err != nil {
				log.Print(err)
				return
			}
			url := r.FormValue("url")
			log.Print("received ", url)
			urls <- url
		}
	}
	stopSong := func(w http.ResponseWriter, r *http.Request) {
		for _, cancel := range cancels {
			cancel()
		}
		cancels = make([]context.CancelFunc, 0)
	}
	http.HandleFunc("/stop", stopSong)
	http.HandleFunc("/play", playSong)
	http.HandleFunc("/add", addSong)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
