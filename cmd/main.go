package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leandro-machado-costa/tl/internal/app/handlers"
	"github.com/leandro-machado-costa/tl/internal/config"
	"github.com/leandro-machado-costa/tl/internal/config/db"
)

// Add the import statement for the router package

func main() {

	err := config.Load()
	if err != nil {
		panic(err)
	}
	err = db.InitDB()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := chi.NewRouter()
	r.Post("/users", handlers.CreateUser)
	r.Get("/users", handlers.GetUsers)
	r.Get(`/users/{id}`, handlers.GetUserByID)
	r.Put(`/users/{id}`, handlers.UpdateUserByID)
	r.Delete(`/users/{id}`, handlers.DeleteUserByID)
	r.Post("/courses", handlers.CreateCourse)
	r.Get("/courses", handlers.GetCourses)
	r.Get(`/courses/{id}`, handlers.GetCourseByID)
	r.Put(`/courses/{id}`, handlers.UpdateCoursesByID)
	r.Delete(`/courses/{id}`, handlers.DeleteCourseByID)

	http.ListenAndServe(fmt.Sprintf(":%s", config.GetServerPort()), r)
}
