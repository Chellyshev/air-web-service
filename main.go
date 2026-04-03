package main

import (
	"fmt"
	"net/http"
	"web/pages"
)

func main() {
	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/styles/", fs)
	http.Handle("/scripts/", fs)

	http.HandleFunc("/", pages.MainPage)
	fmt.Println("Starting server at :8081")
	http.ListenAndServe(":8081", nil)

}
