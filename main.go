package main

import (
	"fmt"
	"net/http"
	"web/pages"
)

func main() {
	fs := http.FileServer(http.Dir("templates"))

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/", pages.AdminMain)
	adminMux.HandleFunc("/panic", pages.PanicPage)

	// set middleware
	adminHandler := pages.AdminAuthMiddleware(adminMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/admin/", adminHandler)
	siteMux.HandleFunc("/", pages.MainPage)
	siteMux.HandleFunc("/login", pages.LoginForm)
	siteMux.HandleFunc("/check", pages.CheckLogin)
	siteMux.HandleFunc("/logout", pages.Logout)
	siteMux.Handle("/styles/", fs)
	siteMux.Handle("/scripts/", fs)
	fmt.Println("Starting server at :8081")
	http.ListenAndServe(":8081", siteMux)

}
