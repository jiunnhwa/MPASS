package main

import (
	"MyLog-M/driver/sqlite"
	delivery "MyLog-M/internal/delivery/http"
	"MyLog-M/internal/repository"
	"MyLog-M/internal/service"
	"log"
	"net/http"
)

func main() {
	db := sqlite.Open("log.db")
	defer sqlite.Close(db)

	repo := repository.NewLogRepo(db)
	s := service.New(*repo)
	handler := delivery.New(s)

	//VIEWS
	http.HandleFunc("/", delivery.Home)

	http.HandleFunc("/api/log", handler.MyLog)
	http.HandleFunc("/api/tail", handler.MyTail)
	http.HandleFunc("/api/view", handler.MyView)

	log.Println(http.ListenAndServe(":558", nil))
}
