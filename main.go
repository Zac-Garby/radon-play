package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Zac-Garby/radon-play/lib"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	http.HandleFunc("/run", handleRun)

	fmt.Println("listening on :3000")

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Println(err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func handleRun(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	if err := lib.HandleConnection(conn); err != nil {
		log.Println(err)
	}
}
