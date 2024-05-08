package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Class struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Time    string `json:"time"`
	Members int    `json:"members"`
}

var users []User
var classes []Class

func main() {
	router := mux.NewRouter()
	
	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/user", createUser).Methods("POST")
	router.HandleFunc("/api/classes", getClasses).Methods("GET")
	router.HandleFunc("/api/class", createClass).Methods("POST")
	router.HandleFunc("/api/bookclass", bookClass).Methods("POST")
	router.HandleFunc("/api/renewmembership", renewMembership).Methods("POST")

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), router); err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	users = append(users, user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func getClasses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(classes)
}

func createClass(w http.ResponseWriter, r *http.Request) {
	var class Class
	_ = json.NewDecoder(r.Body).Decode(&class)
	classes = append(classes, class)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(class)
}

func bookClass(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func renewMembership(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}