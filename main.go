package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/config"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/database"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	siteConfig := &config.SiteConfig{
		DbConnection: dbConn,
		DbQueries:    database.New(dbConn),
	}

	handlerConfig := &handlers.HandlerSiteConfig{SiteConfig: siteConfig}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/create_user", handlerConfig.CreateUserHandler)
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
