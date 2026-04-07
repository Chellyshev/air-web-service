package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"web/pages"

	_ "github.com/lib/pq"
)

func main() {
	pages.Scheduler()
	var err error
	pages.Db, err = sql.Open("postgres", "postgres://air_user:123456@localhost/air_quality?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = pages.Db.Ping()
	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.Dir("templates"))

	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/", pages.AdminMain)

	// set middleware
	adminHandler := pages.AdminAuthMiddleware(adminMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/admin/", adminHandler)
	siteMux.HandleFunc("/", pages.MainPage)
	siteMux.HandleFunc("/login", pages.LoginForm)
	siteMux.HandleFunc("/check", pages.CheckLogin)
	siteMux.HandleFunc("/logout", pages.Logout)
	siteMux.HandleFunc("/charts", pages.Charts)
	siteMux.HandleFunc("/rules", pages.Documents)
	siteMux.HandleFunc("/download/", pages.DownloadGN)
	siteMux.HandleFunc("/api/week", pages.GetWeekData)
	siteMux.Handle("/styles/", fs)
	siteMux.Handle("/scripts/", fs)
	fmt.Println("Starting server at :8081")
	http.ListenAndServe(":8081", siteMux)

}
