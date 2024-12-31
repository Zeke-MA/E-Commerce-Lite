package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	port := os.Getenv("PORT")

	godotenv.Load()

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

}
