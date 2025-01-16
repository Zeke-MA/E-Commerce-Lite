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
	"github.com/Zeke-MA/E-Commerce-Lite/internal/middleware"
	"github.com/gorilla/mux"
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
	middlewareConfig := &middleware.MiddlewareSiteConfig{SiteConfig: siteConfig}

	r := mux.NewRouter()

	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middlewareConfig.CheckUserValidated)

	adminRouter.HandleFunc("/admin/products/add", handlerConfig.AddProduct).Methods("POST")

	r.HandleFunc("/api/create_user", handlerConfig.CreateUserHandler).Methods("POST")
	r.HandleFunc("/api/login", handlerConfig.LoginUserHandler).Methods("POST")
	r.HandleFunc("/api/refresh", handlerConfig.RefreshAccessToken).Methods("POST")
	r.HandleFunc("/api/revoke", handlerConfig.RevokeRefreshToken).Methods("POST")

	r.HandleFunc("/admin/products/remove/{product_id}", handlerConfig.RemoveProduct).Methods("POST")
	r.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	server := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
