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

var secret = []byte("it took the night to believe")

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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if r.Method == "POST" {
		reqBody, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply course information in JSON format"))
			return
		}

		var authenticationInfo map[string]interface{}
		err = json.Unmarshal(reqBody, &authenticationInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}
		isPassenger := authenticationInfo["isPassenger"].(bool)
		email := authenticationInfo["email"].(string)
		var url string
		if isPassenger {
			url = os.Getenv("PASSENGER_MS_HOST") + "/api/v1/passengers?email=" + email
		} else if !isPassenger {
			url = os.Getenv("DRIVER_MS_HOST") + "/api/v1/drivers?email=" + email
		}

		var id int
		var salt string
		var passHash string
		resp, err := http.Get(url)
		if err == nil {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
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

		// Convert salt from hex string to byte array
		decodedSalt, err := hex.DecodeString(salt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}
		// Type assert password into string
		password := authenticationInfo["password"].(string)
		saltedPassword := append([]byte(password), decodedSalt...)

		hashedInput := sha256.Sum256(saltedPassword)
		if fmt.Sprintf("%x", hashedInput) != passHash {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - Authentication failed"))
			return
		}
		token := genJWT(id, email, isPassenger)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{"token": token, "isPassenger": isPassenger}
		// Encode map to json string
		jsonResp, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}
		w.Write(jsonResp)
	}
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	// set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	if r.Method == "POST" {
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

		headerToken := authorizeInfo["authorization"].(string)
		// Decode the jwt and ensure it's readable
		token, err := jwt.Parse(headerToken[7:], func(token *jwt.Token) (interface{}, error) {
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

func main() {
	if os.Getenv("ENVIRONMENT") != "production" {
		os.Setenv("DRIVER_MS_HOST", "http://localhost:5001")
		os.Setenv("PASSENGER_MS_HOST", "http://localhost:5000")
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/login", loginHandler)
	router.HandleFunc("/api/v1/authorize", authHandler)
	fmt.Println("Authentication Microservice")
	fmt.Println("Listening at port 5003")
	log.Fatal(http.ListenAndServe(":5003", router))
}
