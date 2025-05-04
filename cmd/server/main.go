// cmd/server/main.go
package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golangforum/internal/handler"
	"golangforum/internal/repository/postgres"
	"golangforum/internal/usecase"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("ошибка загрузки .env")
	}

	// Подключаемся к базе данных
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")

	// Инициализация репозиториев и use case
	authRepo := postgres.NewRepository(db)
	authUC := usecase.NewAuthUseCase(authRepo)
	authHandler := handler.NewAuthHandler(authUC)

	wsManager := handler.NewWebSocketManager()

	chatRepo := postgres.NewChatRepository(db)
	chatUC := usecase.NewChatUseCase(chatRepo)
	chatHandler := handler.NewChatHandler(chatUC, wsManager)

	topicRepo := postgres.NewTopicRepository(db)
	topicUC := usecase.NewTopicUseCase(topicRepo)
	topicHandler := handler.NewTopicHandler(topicUC)

	postRepo := postgres.NewPostRepository(db)
	postUC := usecase.NewPostUseCase(postRepo)
	postHandler := handler.NewPostHandler(postUC)

	commentRepo := postgres.NewCommentRepository(db)
	commentUC := usecase.NewCommentUseCase(commentRepo)
	commentHandler := handler.NewCommentHandler(commentUC)

	// Настройка маршрутов
	mux := http.NewServeMux()
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)
	mux.HandleFunc("/refresh", authHandler.Refresh)
	mux.HandleFunc("/protected", authHandler.Protected)
	mux.HandleFunc("/chat", chatHandler.ServeWS)
	mux.HandleFunc("/chat/messages", chatHandler.GetAllMessages) // Эндпоинт для получения сообщений

	mux.HandleFunc("/topics", topicHandler.GetAll)
	mux.HandleFunc("/topics/create", topicHandler.Create)
	mux.HandleFunc("/posts", postHandler.GetByTopic)
	mux.HandleFunc("/posts/create", postHandler.Create)
	mux.HandleFunc("/comments", commentHandler.GetByPost)
	mux.HandleFunc("/comments/create", commentHandler.Create)

	// Запуск сервера
	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
