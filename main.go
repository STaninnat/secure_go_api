package main

import (
	"database/sql"
	"embed"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/STaninnat/capstone_project/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB            *database.Queries
	JWTSecret     string
	RefreshSecret string
}

//go:embed static/*
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

	apicfg := apiConfig{}

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

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "static/index.html")
		} else {
			fileServer := http.FileServer(http.FS(staticFiles))
			http.StripPrefix("/", fileServer).ServeHTTP(w, r)
		}
	})

	v1Router := chi.NewRouter()
	if apicfg.DB != nil {
		v1Router.Post("/users", apicfg.handlerUsersCreate)
		v1Router.Get("/users", apicfg.middlewareAuth(apicfg.handlerUsersGet))

		v1Router.Post("/login", apicfg.handlerLogin)

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
