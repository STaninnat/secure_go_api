package main

import (
	"database/sql"
	"embed"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB            *database.Queries
	DBConn        *sql.DB
	JWTSecret     string
	RefreshSecret string
}

//go:embed index.html static/*
var staticFiles embed.FS

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("warning: assuming default configuration. .env unreadable | %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("warning: PORT environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("warning: JWT_SECRET environment variable is not set")
	}

	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		log.Fatal("warning: REFRESH_SECRET environment variable is not set")
	}

	apicfg := apiConfig{
		JWTSecret:     jwtSecret,
		RefreshSecret: refreshSecret,
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("warning: DATABASE_URL environment variable is not set")
		log.Println("Running without CRUD endpoints")
	} else {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("warning: can't connect to database: %v", err)
		}
		dbQueries := database.New(db)
		apicfg.DB = dbQueries
		apicfg.DBConn = db
		log.Println("Connected to database!")
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		f, err := staticFiles.Open("index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "text/html")
		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	router.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		filepath := r.URL.Path[len("/static/"):]

		f, err := staticFiles.Open("static/" + filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		defer f.Close()

		if strings.HasSuffix(filepath, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(filepath, ".js") {
			w.Header().Set("Content-Type", "application/javascript")
		} else if strings.HasSuffix(filepath, ".html") {
			w.Header().Set("Content-Type", "text/html")
		} else if strings.HasSuffix(filepath, ".json") {
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.Header().Set("Content-Type", "text/plain")
		}

		if _, err := io.Copy(w, f); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	v1Router := chi.NewRouter()
	if apicfg.DB != nil {
		v1Router.Post("/users", apicfg.handlerUsersCreate)
		v1Router.Get("/users", apicfg.middlewareAuth(apicfg.handlerUsersGet))

		v1Router.Post("/login", apicfg.handlerLogin)
		v1Router.Post("/logout", apicfg.middlewareAuth(apicfg.handlerLogout))

		v1Router.Post("/refresh", apicfg.handlerRefreshKey)

		v1Router.Post("/posts", apicfg.middlewareAuth(apicfg.handlerPostsCreate))
		v1Router.Get("/posts", apicfg.middlewareAuth(apicfg.handlerPostsGet))
	}

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerError)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
