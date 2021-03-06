package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

//Secret used for generating and decoding JWT
var secret = []byte("it took the night to believe")

// Generates JWT with ID, email, and isPassenger (no expiry)
func genJWT(id int, email string, isPassenger bool) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":          id,
		"email":       email,
		"isPassenger": isPassenger,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}
	return tokenString
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	if r.Method == "POST" {
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply course information in JSON format"))
			return
		}

		// Get POSTed info from body
		var authenticationInfo map[string]interface{}
		err = json.Unmarshal(reqBody, &authenticationInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}
		// Get isPassenger and email from body and cast into proper type
		isPassenger := authenticationInfo["isPassenger"].(bool)
		email := authenticationInfo["email"].(string)
		var url string
		if isPassenger {
			url = os.Getenv("PASSENGER_MS_HOST") + "/api/v1/passengers?email=" + email
		} else if !isPassenger {
			url = os.Getenv("DRIVER_MS_HOST") + "/api/v1/drivers?email=" + email
		}

		// Obtain data for login from passenger/driver microservice
		var id int
		var salt string
		var passHash string
		resp, err := http.Get(url)
		if err == nil {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if len(result) == 0 {
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("403 - Authentication failed"))
					return
				}
				id = int(result["ID"].(float64))
				salt = result["Salt"].(string)
				passHash = result["Password"].(string)
			}
		} else if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Endpoint Unavailable"))
			log.Print(err)
			return
		}

		// Get password from body and type assert to string
		password := authenticationInfo["password"].(string)

		// Convert salt from hex string to byte array
		decodedSalt, err := hex.DecodeString(salt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}
		// Add salt to user-inputted password
		saltedPassword := append([]byte(password), decodedSalt...)

		// Hash user-inputted password and ensure it's the same as the hash in the database
		hashedInput := sha256.Sum256(saltedPassword)
		if fmt.Sprintf("%x", hashedInput) != passHash {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Authentication failed"))
			return
		}
		// Generate JWT from details
		token := genJWT(id, email, isPassenger)
		// Create a cookie containing the JWT
		cookie := &http.Cookie{Name: "jwt", Value: token, MaxAge: 500000, Path: "/", SameSite: http.SameSiteStrictMode, HttpOnly: true}
		// Set cookie in response header
		http.SetCookie(w, cookie)
		w.WriteHeader(200)
		return
	}
}

// For authorization
func authHandler(w http.ResponseWriter, r *http.Request) {
	// set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	if r.Method == "POST" {
		// Other microservices will POST jwt data passed to them to this endpoint
		// to ensure that they are authorized  to perform the action
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply course information in JSON format"))
			return
		}

		var authorizeInfo map[string]interface{}
		err = json.Unmarshal(reqBody, &authorizeInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}

		// Get the JWT from the map
		accessToken := authorizeInfo["authorization"].(string)
		// Decode the jwt and ensure it's readable
		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Invalid Token"))
			return
		}

		// Return ID and isPassenger if decoding is OK
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := int(claims["id"].(float64))
			isPassenger := claims["isPassenger"].(bool)
			json.NewEncoder(w).Encode(map[string]interface{}{"ID": id, "isPassenger": isPassenger})
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Invalid Token"))
			return
		}
	}

	if r.Method == "GET" {
		// This method is solely for frontend to check if user is logged in
		accessToken, err := r.Cookie("jwt")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - No authorized cookie"))
			return
		}
		// Pass token token from cookie and ensure that the jwt is readable
		token, err := jwt.Parse(accessToken.Value, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return secret, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Invalid Token"))
			return
		}

		// Return ID and isPassenger if decoding is OK
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id := int(claims["id"].(float64))
			isPassenger := claims["isPassenger"].(bool)
			json.NewEncoder(w).Encode(map[string]interface{}{"ID": id, "isPassenger": isPassenger})
		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Invalid Token"))
			return
		}
	}
}
func signOutHandler(w http.ResponseWriter, r *http.Request) {
	// set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
	if r.Method == "POST" {
		// Sets response header to delete cookie from browser
		// For signing out of the website
		cookie := &http.Cookie{Name: "jwt", Value: "", MaxAge: 0, Path: "/", SameSite: http.SameSiteStrictMode, HttpOnly: true}
		http.SetCookie(w, cookie)
		w.WriteHeader(200)
	}
}

func main() {
	if os.Getenv("ENVIRONMENT") != "production" {
		os.Setenv("DRIVER_MS_HOST", "http://localhost:5001")
		os.Setenv("PASSENGER_MS_HOST", "http://localhost:5000")
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/login", loginHandler)
	router.HandleFunc("/api/v1/authorize", authHandler).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
	router.HandleFunc("/api/v1/sign-out", signOutHandler).Methods(http.MethodPost, http.MethodOptions)
	router.Use(mux.CORSMethodMiddleware(router))
	fmt.Println("Authentication Microservice")
	fmt.Println("Listening at port 5003")
	log.Fatal(http.ListenAndServe(":5003", router))
}
