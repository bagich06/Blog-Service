package main

import (
	"blog/internal/api"
	"blog/internal/repository"
	"github.com/gorilla/mux"
	"log"
)

func main() {
	db, err := repository.New("postgres://postgres:postgres@localhost:5432/blog_db")
	if err != nil {
		log.Fatal(err)
	}
	api := api.New(mux.NewRouter(), db)
	api.Handle()
	log.Fatal(api.ListenAndServe("localhost:8080"))
}
