package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/leandro-machado-costa/tl/internal/app/handlers"
	"github.com/leandro-machado-costa/tl/internal/configenv"
	"github.com/leandro-machado-costa/tl/internal/configenv/db"

	"github.com/rs/cors"
	"github.com/shaj13/libcache"
	_ "github.com/shaj13/libcache/fifo"

	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/basic"
	"github.com/shaj13/go-guardian/v2/auth/strategies/jwt"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
)

var strategy union.Union
var keeper jwt.SecretsKeeper

// Usage:
// curl  -k http://localhost:9000/user/1 -u admin:admin
// curl  -k http://localhost:9000/auth/token -u admin:admin <obtain a token>
// curl  -k http://localhost:9000/user/1 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsImtpZCI6InNlY3JldC1pZCIsInR5cCI6IkpXVCJ9.eyJFeHRlbnNpb25zIjpudWxsLCJHcm91cHMiOm51bGwsIklEIjoiMSIsIk5hbWUiOiJhZG1pbiIsImF1ZCI6WyIiXSwiZXhwIjoxNzEyNjI0MzIzLCJpYXQiOjE3MTI2MjQwMjMsIm5iZiI6MTcxMjYyNDAyMywic3ViIjoiMSJ9.t9DC-lcqVvQeET6xLcZGZgHt1rjDEUjhNwk1g7OBeB0 "

func main() {

	setupGoGuardian()

	err := configenv.Load()
	if err != nil {
		panic(err)
	}
	err = db.InitDB()

	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := chi.NewRouter()

	r.Use(cors.Default().Handler)

	r.Get("/auth/token", middleware(http.HandlerFunc(createToken)))

	r.Get("/users", handlers.GetUsers)
	r.Get(`/users/{id}`, handlers.GetUserByID)
	r.Put(`/users/{id}`, middleware(http.HandlerFunc(handlers.UpdateUserByID)))
	r.Post("/users", middleware(http.HandlerFunc(handlers.CreateUser)))
	r.Delete(`/users/{id}`, middleware(http.HandlerFunc(handlers.DeleteUserByID)))
	r.Get("/courses", handlers.GetCourses)
	r.Get(`/courses/{id}`, handlers.GetCourseByID)
	r.Post("/courses", middleware(http.HandlerFunc(handlers.CreateCourse)))
	r.Put(`/courses/{id}`, middleware(http.HandlerFunc(handlers.UpdateCoursesByID)))
	r.Delete(`/courses/{id}`, middleware(http.HandlerFunc(handlers.DeleteCourseByID)))

	http.ListenAndServe(fmt.Sprintf(":%s", configenv.GetServerPort()), r)
}

func createToken(w http.ResponseWriter, r *http.Request) {
	u := auth.User(r)
	token, _ := jwt.IssueAccessToken(u, keeper)
	body := fmt.Sprintf("token: %s \n", token)
	w.Write([]byte(body))
}

func setupGoGuardian() {
	keeper = jwt.StaticSecret{
		ID:        "secret-id",
		Secret:    []byte("secret"),
		Algorithm: jwt.HS256,
	}
	cache := libcache.FIFO.New(0)
	cache.SetTTL(time.Minute * 5)
	basicStrategy := basic.NewCached(validateUser, cache)
	jwtStrategy := jwt.New(cache, keeper)
	strategy = union.New(jwtStrategy, basicStrategy)
}

func validateUser(ctx context.Context, r *http.Request, userName, password string) (auth.Info, error) {
	// here connect to db or any other service to fetch user and validate it.
	if userName == "admin" && password == "admin" {
		return auth.NewDefaultUser("admin", "1", nil, nil), nil
	}

	return nil, fmt.Errorf("Invalid credentials")
}

func middleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing Auth Middleware")
		_, user, err := strategy.AuthenticateRequest(r)
		if err != nil {
			fmt.Println(err)
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}
		log.Printf("User %s Authenticated\n", user.GetUserName())
		r = auth.RequestWithUser(user, r)
		next.ServeHTTP(w, r)
	})
}
