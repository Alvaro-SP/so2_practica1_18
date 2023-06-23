package main

import (
	"Backend/routes"
	"log"
	"net/http"
	"time"
  "fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()
	routes.ComandoRoute(router) //add this
	//CORS

	corsWrapper := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*"},
	})

	srv := &http.Server{
		Handler:      corsWrapper.Handler(router),
		Addr:         ":4000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//enableCORS(router)
  fmt.Println("Server on port 4000") 
	log.Fatal(srv.ListenAndServe())
}
