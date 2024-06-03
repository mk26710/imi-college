package main

import (
	"imi/college/internal/enum"
	"imi/college/internal/env"
	"imi/college/internal/handlers"
	mw "imi/college/internal/middleware"
	"imi/college/internal/models"
	"log"
	"net/http"

	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(env.DSN()), &gorm.Config{TranslateError: true})
	if err != nil {
		log.Fatalln("Couldn't connect to postgres database")
	}

	models.AutoMigrate(db)

	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.CleanPath)
	r.Use(chimw.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5173/"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	h := handlers.Create(db)

	// Public routes group
	r.Group(func(r chi.Router) {
		r.Post("/users", handlers.APIHandler(h.Users.Create))
		r.Post("/tokens", handlers.APIHandler(h.Tokens.Create))
		r.Delete("/tokens", handlers.APIHandler(h.Tokens.Delete))

		r.Group(func(r chi.Router) {
			r.Get("/dictionaries/regions", handlers.APIHandler(h.Dictionaries.ReadRegions))
			r.Get("/dictionaries/towntypes", handlers.APIHandler(h.Dictionaries.ReadTownTypes))
		})
	})

	// Authentication required
	r.Group(func(r chi.Router) {
		r.Use(mw.RequireUser(db))

		r.Get("/users/@me", handlers.APIHandler(h.Users.ReadMe))
		r.Get("/users/@me/address", handlers.APIHandler(h.Address.ReadMe))

		r.With(mw.RequirePermissions(enum.PermissionViewUser)).Get("/users/{id}", handlers.APIHandler(h.Users.Read))

		r.Post("/files", handlers.APIHandler(h.Files.CreateFile))

		r.Post("/documents/identity", handlers.APIHandler(h.Documents.CreateDocumentIdentity))
	})

	srv := http.Server{
		Addr:    env.Addr(),
		Handler: r,
	}

	log.Printf("Lisetning on http://%s\n", env.Addr())
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
