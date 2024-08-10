package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mfauzidr/coffeeshop-go-backend/internal/routers"
	"github.com/mfauzidr/coffeeshop-go-backend/pkg"
)

func main() {
	log.Println("Starting application...")

	db, err := pkg.Posql()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	log.Println("Database initialized successfully")

	router := routers.New(db)
	server := pkg.Server(router)

	log.Println("Server is running at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
