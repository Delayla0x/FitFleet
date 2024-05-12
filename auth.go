package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/joho/godotenv"
    "github.com/gorilla/mux"
)

var jwtKey []byte

type UserClaims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

func init() {
    if err := godotenv.Load(); err != nil {
        log.Printf("No .env file found, running with defaults")
    }
    jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/login", LoginHandler).Methods("POST")
    r.HandleFunc("/dashboard", TokenVerifyMiddleWare(DashboardHandler)).Methods("GET")

    log.Fatal(http.ListenAndServe(":8080", r))
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
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(tokenString))
    } else {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
    claims, ok := r.Context().Value("claims").(*UserClaims)
    if !ok {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if claims.Role == "staff" {
        w.Write([]byte(fmt.Sprintf("Hello, %s! Welcome to the staff dashboard.", claims.Username)))
    } else if claims.Role == "regular" {
        w.Write([]byte(fmt.Sprintf("Hello, %s! Welcome to your dashboard.", claims.Username)))
    } else {
        w.WriteHeader(http.StatusForbidden)
        return
    }
}

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")

        claims := &UserClaims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        
        ctx := r.Context()
        ctx = context.WithValue(ctx, "claims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}