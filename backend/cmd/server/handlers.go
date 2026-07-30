package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/brain-trainer/backend/internal/db"
	"github.com/brain-trainer/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct{
	Email 					string `json:"email"`
	Password 		  string `json:"password"`
	Username 		 string `json:"username,omitempty"`
}

type AuthResponse struct{
	Token 				string `json:"token"`
	User 			   	 models.User `json:"user"`
}

type ScoreSubmission struct{
	GameType	 		  string 				`json:"gameType"`
	Score 						 int 					   `json:"score"`
	Level 						  int 						`json:"level"`
	MaxStreak 			   int 						`json:"maxStreak"`
	TotalAttempts		int 					 `json:"totalAttempts"`
	CorrectAttempts  int 					  `json:"correctAttempts"`
	SessionData 		json.RawMessage 	`json:"sessionData,omitempty"`
}

// Authentication handlers 

func(s *Server) handleRegister(w http.ResponseWriter, r *http.Request){
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// validate input 
	if req.Email == "" || req.Password =="" || req.Username ==""{
		writeError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	if len(req.Password) <8{
		writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil{
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	// Create user 
	user, err := s.queries.CreateUser(r.Context(), db.CreateUserParams{
		Username: 		req.Username,
		Email:				   req.Email,
		PasswordHash:  string(hashedPassword),
	})
}


func writeJSON(w http.ResponseWriter, status int, data interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}


func writeError(w http.ResponseWriter, status int, message string){
	writeJSON(w, status, map[string]string{"error":message})
}