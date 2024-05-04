package main

import (
	mw "imi/college/internal/middleware"
	"imi/college/internal/models"
	"imi/college/internal/routes"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RoutesHandlers struct {
	User    *routes.UserHandler
	Session *routes.SessionHandler
	File    *routes.FilesHandler
}

func CreateHandlers(db *gorm.DB) RoutesHandlers {
	return RoutesHandlers{
		User:    routes.NewUserHandler(db),
		Session: routes.NewSessionHandler(db),
		File:    routes.NewFilesHandler(db),
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DB_DSN")), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalln("Couldn't connect to postgres database")
	}

	db.AutoMigrate(&models.User{}, &models.Password{}, &models.UserSession{})

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.CleanPath)
	r.Use(chimw.Recoverer)

	h := CreateHandlers(db)

	// Public routes group
	r.Group(func(r chi.Router) {
		r.Post("/users", h.User.CreateUser)
		r.Post("/session", h.Session.CreateSession)
	})

	// Authentication required
	r.Group(func(r chi.Router) {
		r.Use(mw.EnsureUserSession(db))
		r.Get("/users/{id}", h.User.ReadUser)
		r.Post("/upload", h.File.CreateFile)
	})

	srv := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: r,
	}

	log.Println("Lisetning on 8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
