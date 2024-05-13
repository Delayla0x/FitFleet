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

type FitnessClass struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Time    string `json:"time"`
	Members int    `json:"members"`
}

var registeredUsers []User
var availableClasses []FitnessClass

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/users", GetAllUsers).Methods("GET")
	router.HandleFunc("/api/user", RegisterNewUser).Methods("POST")
	router.HandleFunc("/api/classes", GetAllClasses).Methods("GET")
	router.HandleFunc("/api/class", CreateNewClass).Methods("POST")
	router.HandleFunc("/api/bookclass", BookAClass).Methods("POST")
	router.HandleFunc("/api/renewmembership", RenewMembership).Methods("POST")

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(registeredUsers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	registeredUsers = append(registeredUsers, newUser)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetAllClasses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(availableClasses); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateNewClass(w http.ResponseWriter, r *http.Request) {
	var newClass FitnessClass
	if err := json.NewDecoder(r.Body).Decode(&newClass); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	availableClasses = append(availableClasses, newClass)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(newClass); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func BookAClass(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func RenewMembership(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}