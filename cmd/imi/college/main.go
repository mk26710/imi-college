package main

import (
	"imi/college/internal/env"
	"imi/college/internal/handlers"
	"imi/college/internal/httpx"
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
		r.Post("/users", httpx.APIHandler(h.Users.Create))
		r.Post("/tokens", httpx.APIHandler(h.Tokens.Create))
		r.Delete("/tokens", httpx.APIHandler(h.Tokens.Delete))

		r.Group(func(r chi.Router) {
			r.Get("/dictionaries/regions", httpx.APIHandler(h.Dictionaries.ReadRegions))
			r.Get("/dictionaries/towntypes", httpx.APIHandler(h.Dictionaries.ReadTownTypes))
			r.Get("/dictionaries/genders", httpx.APIHandler(h.Dictionaries.ReadGenders))
		})
	})

	// Authentication required
	r.Group(func(r chi.Router) {
		r.Use(mw.RequireUser(db))

		r.Get("/users/{id}", httpx.APIHandler(h.Users.Read))
		r.Get("/users/{id}/address", httpx.APIHandler(h.Address.Read))

		r.Get("/users/{userId}/documents/identity", httpx.APIHandler(h.Identities.Read))
		r.Post("/users/{userId}/documents/identity", httpx.APIHandler(h.Identities.Create))

		r.Put("/users/{id}/address", httpx.APIHandler(h.Address.CreateOrUpdate))
		r.Put("/users/{id}/details", httpx.APIHandler(h.Users.CreateOrUpdate))

		r.Post("/files", httpx.APIHandler(h.Files.CreateFile))

		r.Post("/documents/identity", httpx.APIHandler(h.Documents.CreateDocumentIdentity))
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
