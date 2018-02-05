package controllers

import (
	"database/sql"
	"github.com/golang/standard-rest-api/utils/caching"
	"net/http"
	"encoding/json"
	"github.com/golang/standard-rest-api/requests"
	"github.com/golang/standard-rest-api/repositories"
	"log"
	"github.com/golang/standard-rest-api/utils/crypto"
	"time"
	"fmt"
	"strconv"
)

type UserController struct {
	DB *sql.DB
	Cache caching.Cache
}

func NewUserController(db *sql.DB, c caching.Cache) *UserController {
	return &UserController{
		DB: db,
		Cache: c,
	}
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var rr requests.RegisterRequest
	err := decoder.Decode(&rr)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := repositories.CreateUser(uc.DB, rr.Email, rr.Name, rr.Password)
	if err != nil {
		log.Fatalf("Add user to database error:%s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	token, err := crypto.GenerateToken()
	if err != nil {
		log.Fatalf("Generate token error:%s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	oneMonth := time.Duration(60*60*24*30) * time.Second
	err = uc.Cache.Set(fmt.Sprintf("token_%s", token), strconv.Itoa(id), oneMonth)
	if err != nil {
		log.Fatalf("Add token to redis error:%s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	p := map[string]string {
		"Token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var lr requests.LoginRequest
	err := decoder.Decode(lr)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := repositories.GetPrivateUserDetailByEmail(uc.DB, lr.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusBadRequest)
			return
		}
		log.Fatal("Create user error:%s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	password := crypto.HashPassword(lr.Password, user.Salt)
	if user.Password != password {
		http.Error(w, "Invalid username or password", http.StatusBadRequest)
		return
	}

	token, err := crypto.GenerateToken()
	if err != nil {
		log.Fatal("Create user error:%s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	oneMonth := time.Duration(60*60*24*30) * time.Second
	err = uc.Cache.Set(fmt.Sprintf("token_%s", token), strconv.Itoa(user.ID), oneMonth)
	if err != nil {
		log.Fatalf("Create user error:%s", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	p := map[string]string {
		"token": token,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
