package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alvar0liveira/rssagregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Can not connect to database")
	}

	db := database.New(conn)

	apiConfig := apiConfig{
		DB: db,
	}

	go startScraper(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/users", apiConfig.handleCreateUser)
	v1Router.Get("/users", apiConfig.middlewareAuth(apiConfig.handleGetUserByApiKey))

	v1Router.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handleCreateFeed))
	v1Router.Get("/feeds", apiConfig.handleGetFeeds)
	v1Router.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.handleCreateFeedFollow))
	v1Router.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.handleGetFeedFollows))

	v1Router.Get("/posts", apiConfig.middlewareAuth(apiConfig.handleGetPostsForUser))

	v1Router.Delete("/feed_follows/{feedFollowId}", apiConfig.middlewareAuth(apiConfig.handleDeleteFeedFollows))
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerErr)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Port:", portString)
}
