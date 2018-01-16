package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Zac-Garby/radon-play/lib"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/run/", handleRun)

	fmt.Printf("listening on :%s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	jobQuery, ok := q["job"]
	if !ok {
		jobQuery = []string{"exec"}
	}

	job := jobQuery[0]

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if err := lib.HandleConnection(conn, job); err != nil {
		log.Println(err)
	}
}
