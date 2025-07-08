package routes

import (
	"github.com/SHIVAM-GOUR/gbt-master-backend/handlers"
	"github.com/SHIVAM-GOUR/gbt-master-backend/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORSMiddleware())

	// Special API
	r.Route("/riya", func(r chi.Router) {
		r.Get("/", handlers.HiRiya)
	})

	// Class routes
	r.Route("/classes", func(r chi.Router) {
		r.Get("/", handlers.GetClasses)
		r.Post("/", handlers.CreateClass)
		r.Get("/{id}", handlers.GetClass)
		r.Put("/{id}", handlers.UpdateClass)
		r.Delete("/{id}", handlers.DeleteClass)
	})

	r.Route("/inquiry", func(r chi.Router) {
		r.Post("/", handlers.CreateInquiry)
	})

	return r
}
