package main

import (
	"imi/college/internal/middleware"
	"imi/college/internal/models"
	"imi/college/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalln("Couldn't connect to postgres database")
	}

	db.AutoMigrate(&models.User{}, &models.Password{}, &models.UserSession{})

	mux := http.NewServeMux()
	stack := middleware.CreateStack(
		middleware.Logging,
	)

	users := routes.NewUserHandler(db)
	sessions := routes.NewSessionHandler(db)

	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("GET /users/{id}", users.ReadUser)

	mux.HandleFunc("POST /session", sessions.CreateSession)

	server := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: stack(mux),
	}

	log.Println("Lisetning on 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
