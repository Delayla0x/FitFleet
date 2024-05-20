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
        log.Printf("Warning: No .env file found. Assuming environment variables are set externally. Error: %v", err)
    }

    secret, present := os.LookupEnv("JWT_SECRET_KEY")
    if !present || secret == "" {
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
            httpError(responseWriter, "Error creating the token", http.StatusInternalServerError)
            log.Printf("Error signing token: %v", err)
            return
        }

        responseWriter.Header().Set("Content-Type", "application/json")
        responseWriter.WriteHeader(http.StatusOK)
        responseWriter.Write([]byte(tokenString))
    } else {
        httpError(responseWriter, "Invalid username or password", http.StatusUnauthorized)
    }
}

func DashboardHandler(responseWriter http.ResponseWriter, request *http.Request) {
    claims, ok := request.Context().Value("userClaims").(*UserClaims)
    if !ok {
        httpError(responseWriter, "Error retrieving user claims", http.StatusInternalServerError)
        return
    }

    var responseMessage string
    switch claims.Role {
    case "staff":
        responseMessage = fmt.Sprintf("Hello, %s! Welcome to the staff dashboard.", claims.Username)
    case "regular":
        responseMessage = fmt.Sprintf("Hello, %s! Welcome to your dashboard.", claims.Username)
    default:
        httpError(responseWriter, "Access denied", http.StatusForbidden)
        return
    }

    responseWriter.Header().Set("Content-Type", "text/plain")
    responseWriter.WriteHeader(http.StatusOK)
    responseWriter.Write([]byte(responseMessage))
}

func TokenVerificationMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(responseWriter http.ResponseWriter, request *http.Request) {
        tokenString := request.Header.Get("Authorization")

        if tokenString == "" {
            httpError(responseWriter, "Authorization token required", http.StatusUnauthorized)
            return
        }

        userClaims := &UserClaims{}

        token, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
            return jwtSecretKey, nil
        })

        if err != nil || !token.Valid {
            httpError(responseWriter, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(request.Context(), "userClaims", userClaims)
        next.ServeHTTP(responseWriter, request.WithContext(ctx))
    }
}

func httpError(writer http.ResponseWriter, message string, statusCode int) {
    writer.WriteHeader(statusCode)
    err := json.NewEncoder(writer).Encode(ErrorResponse{ErrorMessage: message})
    if err != nil {
        log.Printf("Error writing error response: %v", err)
    }
}