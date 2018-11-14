package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	// "github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
)

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ua := req.Header.Get("Content-Type")
	log.Print(ua)
	if !strings.Contains(ua, "application/json") {
		responseCode(w, http.StatusUnsupportedMediaType)
		return
	}

	user := &User{}
	err = json.Unmarshal(data, user)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := User{}
	err = users.Find(bson.M{"username": user.Username}).One(&result)
	if err != nil {
		responseError(w, "User not found", http.StatusNotFound)
		return
	}
	rez := CheckPasswordHash(user.Password, result.Password)
	if !rez {
		// responseError(w, "Wrong password!", http.StatusUnauthorized) 
		responseError(w, "Wrong password!", http.StatusUnauthorized) 
		return
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"password": user.Password,
	})
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
			fmt.Println(error)
	}
	json.NewEncoder(w).Encode(JwtToken{Token: tokenString})
}

func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {
	// params := req.URL.Query()
	// log.Print(req.Header.Get("token"))
	token, _ := jwt.Parse(req.Header.Get("authorization"), func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
			}
			return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var user User
			mapstructure.Decode(claims, &user)
			json.NewEncoder(w).Encode(user)
	} else {
			json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
	}
}

func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			authorizationHeader := req.Header.Get("authorization")
			if authorizationHeader != "" {
					bearerToken := strings.Split(authorizationHeader, " ")
					if len(bearerToken) == 2 {
							token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
									if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
											return nil, fmt.Errorf("There was an error")
									}
									return []byte("secret"), nil
							})
							if error != nil {
									json.NewEncoder(w).Encode(Exception{Message: error.Error()})
									return
							}
							if token.Valid {
									context.Set(req, "decoded", token.Claims)
									next(w, req)
							} else {
									json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
							}
					}
			} else {
					json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
			}
	})
}

func TestEndpoint(w http.ResponseWriter, req *http.Request) {
	decoded := context.Get(req, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)
	json.NewEncoder(w).Encode(user)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}