package main

import (
	"database/sql"
	"log"
	"log/slog"
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
	middlewareConfig := &middleware.MiddlewareSiteConfig{
		SiteConfig: siteConfig,
		Logger:     slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
	}

	r := mux.NewRouter()
	r.Use(middlewareConfig.LogIncomingRequest)

	adminRouter := r.PathPrefix("/admin").Subrouter()

	adminRouter.Use(middlewareConfig.CheckUserValidated)

	adminRouter.HandleFunc("/products/{product_id}/add", handlerConfig.AdminAddProduct).Methods("POST")
	adminRouter.HandleFunc("/products/{product_id}/remove", handlerConfig.AdminRemoveProduct).Methods("DELETE")
	adminRouter.HandleFunc("/products/{product_id}/price", handlerConfig.AdminChangePrice).Methods("PATCH")
	adminRouter.HandleFunc("/products/{product_id}/quantity/add", handlerConfig.AdminAddQuantity).Methods("PATCH")
	adminRouter.HandleFunc("/products/{product_id}/quantity/remove", handlerConfig.AdminRemoveQuantity).Methods("PATCH")

	r.HandleFunc("/api/create_user", handlerConfig.CreateUserHandler).Methods("POST")
	r.HandleFunc("/api/login", handlerConfig.LoginUserHandler).Methods("POST")
	r.HandleFunc("/api/refresh", handlerConfig.RefreshAccessToken).Methods("POST")
	r.HandleFunc("/api/revoke", handlerConfig.RevokeRefreshToken).Methods("POST")

	r.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	server := http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

}
