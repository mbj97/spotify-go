package main

import (
	"fmt"
	"net/http"
	"spotify-go/handler"
)

func main() {
	http.HandleFunc("/auth", handler.AuthHandler)
	http.HandleFunc("/", handler.CallbackHandler)
	http.HandleFunc("/playing", handler.CurrentlyPlayingHandler)
	http.HandleFunc("/start",handler.StartTaskHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
	select {}
}
