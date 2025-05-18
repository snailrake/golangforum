package main

import (
	"database/sql"
	"golangforum/internal/repository/impl"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "golangforum/docs"
	"golangforum/internal/client"
	"golangforum/internal/handler"
	usecaseImpl "golangforum/internal/usecase/impl"
)

var logger zerolog.Logger

func init() {
	logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

// @title API сервиса форума
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	if err := godotenv.Load(); err != nil {
		logger.Fatal().Err(err).Msg("failed to load .env")
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open database connection")
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		logger.Fatal().Err(err).Msg("failed to ping database")
	}

	authClient, err := client.NewClient(os.Getenv("AUTH_SERVICE_ADDR"))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to initialize auth client")
	}

	chatHandler := handler.NewChatHandler(
		usecaseImpl.NewChatUseCase(impl.NewChatRepository(db)),
		authClient,
	)
	topicHandler := handler.NewTopicHandler(
		usecaseImpl.NewTopicUseCase(impl.NewTopicRepository(db)),
	)
	postHandler := handler.NewPostHandler(
		usecaseImpl.NewPostUseCase(impl.NewPostRepository(db)),
		authClient,
	)
	commentHandler := handler.NewCommentHandler(
		usecaseImpl.NewCommentUseCase(impl.NewCommentRepository(db)),
		authClient,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/chat", chatHandler.ServeWS)
	mux.HandleFunc("/chat/messages", chatHandler.GetAllMessages)
	mux.HandleFunc("/topics", topicHandler.GetAll)
	mux.HandleFunc("/topics/create", topicHandler.Create)
	mux.HandleFunc("/topics/delete", topicHandler.Delete)
	mux.HandleFunc("/posts", postHandler.GetByTopic)
	mux.HandleFunc("/posts/all", postHandler.GetAll)
	mux.HandleFunc("/posts/create", postHandler.Create)
	mux.HandleFunc("/posts/delete", postHandler.Delete)
	mux.HandleFunc("/comments", commentHandler.GetByPost)
	mux.HandleFunc("/comments/create", commentHandler.Create)
	mux.HandleFunc("/comments/delete", commentHandler.Delete)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	logger.Info().Msg("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug().Str("method", r.Method).Str("url", r.URL.String()).Msg("handling request")
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
