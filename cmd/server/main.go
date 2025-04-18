package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golangforum/internal/handler"
	"golangforum/internal/repository/postgres"
	"golangforum/internal/usecase"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используется системное окружение")
	}

	postgresURI := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatal(err)
	}
	fmt.Println("Connected to database")

	repo := postgres.NewRepository(db)
	authUC := usecase.NewAuthUseCase(repo)
	authHandler := handler.NewAuthHandler(authUC)

	mux := http.NewServeMux()
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/refresh", authHandler.Refresh)
	mux.HandleFunc("/protected", authHandler.Protected)

	handlerWithCors := corsMiddleware(mux)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handlerWithCors))
}
