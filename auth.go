package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

var jwtKey []byte

type UserClaims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

type ErrorResponse struct {
    Error string `json:"error"`
}

func init() {
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: No .env file found, running with defaults or existing environment variables")
    }
    
    secret := os.Getenv("JWT_SECRET_KEY")
    if secret == "" {
        log.Fatal("JWT_SECRET_KEY is not set. Exiting application.")
    }
    jwtKey = []byte(secret)
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/login", LoginHandler).Methods("POST")
    r.HandleFunc("/dashboard", TokenVerifyMiddleWare(DashboardHandler)).Methods("GET")

    log.Printf("Server is running on port 8080")
    err := http.ListenAndServe(":8080", r)
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var username, password string
    
    username = "user1"
    password = "password"

    if username == "user1" && password == "password" {
        userRole := "regular"

        expirationTime := time.Now().Add(1 * time.Hour)
        claims := &UserClaims{
            Username: username,
            Role:     userRole,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtKey)

        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(ErrorResponse{Error: "Error creating the token"})
            log.Printf("Error signing token: %v", err)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(tokenString))
    } else {
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid username or password"})
        return
    }
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(*UserClaims)
    if !ok {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(ErrorResponse{Error: "Error retrieving claims"})
        return
    }

    var responseMsg string
    if claims.Role == "staff" {
        responseMsg = fmt.Sprintf("Hello, %s! Welcome to the staff dashboard.", claims.Username)
    } else if claims.Role == "regular" {
        responseMsg = fmt.Sprintf("Hello, %s! Welcome to your dashboard.", claims.Username)
    } else {
        w.WriteHeader(http.StatusForbidden)
        json.NewEncoder(w).Encode(ErrorResponse{Error: "Access denied"})
        return
    }

    w.Write([]byte(responseMsg))
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")

        claims := &UserClaims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid token"})
            return
        }

        ctx := r.Context()
        ctx = context.WithValue(ctx, "claims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}