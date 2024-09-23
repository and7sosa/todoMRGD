package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	// create storage
	db := ConnectDB(os.Getenv("DB_NAME"), os.Getenv("DB_URI"))
	// get collection
	s := NewStorage(os.Getenv("COLLECTION_NAME"), db)
	// create router
	apiServer := NewAPIServer(port, s)
	// serve router
	log.Printf("Starting server on port %v...\n", port)
	apiServer.Serve()
}
