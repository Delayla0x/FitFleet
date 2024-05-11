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
	Password string `json:"password"` // Note: Storing passwords as plain text is insecure; consider using hashed passwords.
}

type FitnessClass struct { // Renamed from Class to FitnessClass for clarity
	ID      string `json:"id"`
	Name    string `json:"name"`
	Time    string `json:"time"`
	Members int    `json:"members"`
}

var registeredUsers []User           // Renamed from users to registeredUsers for clarity
var availableClasses []FitnessClass  // Renamed from classes to availableClasses for clarity

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
	json.NewEncoder(w).Encode(registeredUsers)
}

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	_ = json.NewDecoder(r.Body).Decode(&newUser)
	registeredUsers = append(registeredUsers, newUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newUser)
}

func GetAllClasses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(availableClasses)
}

func CreateNewClass(w http.ResponseWriter, r *http.Request) {
	var newClass FitnessClass
	_ = json.NewDecoder(r.Body).Decode(&newClass)
	availableClasses = append(availableClasses, newClass)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newClass)
}

func BookAClass(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func RenewMembership(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}