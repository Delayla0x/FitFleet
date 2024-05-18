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

var jwtSecretKey []byte

type UserClaims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

type ErrorResponse struct {
    ErrorMessage string `json:"error"`
}

func init() {
    if err := godotenv.Load(); err != nil {
        log.Printf("Warning: No .env file found. Running with defaults or existing environment variables.")
    }

    secret := os.Getenv("JWT_SECRET_KEY")
    if secret == "" {
        log.Fatal("JWT_SECRET_KEY is not set. Exiting application.")
    }
    jwtSecretKey = []byte(secret)
}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/login", LoginHandler).Methods("POST")
    router.HandleFunc("/dashboard", TokenVerificationMiddleware(DashboardHandler)).Methods("GET")

    log.Printf("Server is running on port 8080")
    if err := http.ListenAndServe(":8080", router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func LoginHandler(responseWriter http.ResponseWriter, request *http.Request) {
    demoUsername := "user1"
    demoPassword := "password"

    // Mimicking a check against a user store or database.
    // This part should involve checking request body or parameters for actual username/password provided by the user.
    
    if demoUsername == "user1" && demoPassword == "password" {
        userRole := "regular"

        expirationTime := time.Now().Add(1 * time.Hour)
        claims := &UserClaims{
            Username: demoUsername,
            Role:     userRole,
            StandardClaims: jwt.StandardClaims{
                ExpiresAt: expirationTime.Unix(),
            },
        }

        token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
        tokenString, err := token.SignedString(jwtSecretKey)

        if err != nil {
            responseWriter.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(responseWriter).Encode(ErrorResponse{ErrorMessage: "Error creating the token"})
            log.Printf("Error signing token: %v", err)
            return
        }

        responseWriter.Header().Set("Content-Type", "application/json")
        responseWriter.WriteHeader(http.StatusOK)
        responseWriter.Write([]byte(tokenString))
    } else {
        responseWriter.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(responseWriter).Encode(ErrorResponse{ErrorMessage: "Invalid username or password"})
        return
    }
}

func DashboardHandler(responseWriter http.ResponseWriter, request *http.Request) {
    claims, ok := request.Context().Value("userClaims").(*UserClaims)
    if !ok {
        responseWriter.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(responseWriter).Encode(ErrorResponse{ErrorMessage: "Error retrieving user claims"})
        return
    }

    var responseMessage string
    if claims.Role == "staff" {
        responseMessage = fmt.Sprintf("Hello, %s! Welcome to the staff dashboard.", claims.Username)
    } else if claims.Role == "regular" {
        responseMessage = fmt.Sprintf("Hello, %s! Welcome to your dashboard.", claims.Username)
    } else {
        responseWriter.WriteHeader(http.StatusForbidden)
        json.NewEncoder(responseWriter).Encode(ErrorResponse{ErrorMessage: "Access denied"})
        return
    }

    responseWriter.Header().Set("Content-Type", "text/plain")
    responseWriter.WriteHeader(http.StatusOK)
    responseWriter.Write([]byte(responseMessage))
}

func TokenVerificationMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(responseWriter http.ResponseWriter, request *http.Request) {
        tokenString := request.Header.Get("Authorization")

        // Basic check to ensure the token was provided
        if tokenString == "" {
            responseWriter.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(responseWriter).Encode(ErrorResponse{ErrorMessage: "Authorization token required"})
            return
        }

        userClaims := &UserClaims{}

        token, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
            return jwtSecretKey, nil
        })

        if err != nil || !token.Valid {
            responseWriter.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(responseWriter).Encode(ErrorResponse{ErrorMessage: "Invalid token"})
            return
        }

        ctx := context.WithValue(request.Context(), "userClaims", userClaims)
        next.ServeHTTP(responseWriter, request.WithContext(ctx))
    }
}