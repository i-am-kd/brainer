package main 

import(
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/brainer/backend/internal/db"
)

type Server struct{
	db 				*sql.DB
	queries 	*db.Queries
}

func main(){
	//load local environment 
	dbURL := os.Getenv("DB_URL")
	if dbURL ==""{
		dbURL = "postgres://postgres:postgres@localhost:5432/brain_trainer?sslmode=disable"
	}

	//initialize database connection with coneection pooling 
	dbConn, err := sql.Open("pgx", dbURL)
	if err != nil{
		log.Fatal("failed to connect to database : %w", err)
	}
	defer dbConn.Close()

	dbConn.SetMaxOpenConns(25)
	dbConn.SetMaxIdleConns(5)
	dbConn.SetConnMaxIdleTime(2*time.Minute)

	// verify the connection 
	if err :=dbConn.Ping(); err !=nil{
		log.Fatal("failed to ping database: %v", err)
	}

	// initlaize sqlc queries 
	queries := db.new(dbConn)
	// initialize the server 
	server := &Server{
		db: 		dbConn, 
		queries: queries,
	}

	mux := http.NewServeMux()

	// APi routes 
	api := http.NewServeMux()
	api.HandleFunc("POST /auth/register", server.handleRegister)
	api.HandleFunc("POST /auth/login", server.handleLogin)
	api.HandleFunc("POST /auth/logout", server.handleLogout)
	api.HandleFunc("GET /auth/verify", server.handleVerify)
	api.HandleFunc("POST /score/submit", server.handleSubmitScore)
	api.HandleFunc("GET /leaderboard/{gameType}", server.handleGetLeaderboard)
	api.HandleFunc("GET /users/{userId}/stats", server.handleUserStats)
	api.HandleFunc("GET /health", server.handleHealth)

	// apply middleware 
	handler := corsMiddleware(api)
	handler = loggingMiddleware(handler)

	// mout api 
	mux.Handle("/api/", http.StripPrefix("/api/", handler))

	// create http server 
	srv := &http.Server{
		Addr: 				":8080",
		Handler: 		mux,
		ReadTimeout: 15*time.Second,
		WriteTimeout: 15*time.Second,
		IdleTimeout: 60*time.Second,
	}

	// start server in goroutine 
	go func(){
		log.Printf("Server starting on %s", srv.Addr)
		if err := srv.ListenAndServe(); err !=nil && err !=http.ErrServerClosed{
			log.Fatal("Server failed: %v", err)
		}
	}()

	//graceful shutdown 
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Printf("Server shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil{
		log.Fatal("Server shutdown failed: %v", err)
	}
}