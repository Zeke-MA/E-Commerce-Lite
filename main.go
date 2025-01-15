package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

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
		DbConnection:       dbConn,
		DbQueries:          database.New(dbConn),
		RefreshTokenExpiry: time.Hour * 24 * 14,
		JWTExpiry:          time.Hour,
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}

	handlerConfig := &handlers.HandlerSiteConfig{SiteConfig: siteConfig}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/create_user", handlerConfig.CreateUserHandler)
	mux.HandleFunc("POST /api/login", handlerConfig.LoginUserHandler)
	mux.HandleFunc("POST /api/refresh", handlerConfig.RefreshAccessToken)
	mux.HandleFunc("POST /api/revoke", handlerConfig.RevokeRefreshToken)
	mux.HandleFunc("POST /admin/products/add", handlerConfig.AddProduct)
	mux.HandleFunc("POST /admin/products/remove/{product_id}", handlerConfig.RemoveProduct)
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
