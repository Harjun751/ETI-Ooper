package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
			panic(err.Error())
		}
		isPassenger := authenticationInfo["isPassenger"].(bool)
		email := authenticationInfo["email"].(string)
		var url string
		if isPassenger {
			url = "http://localhost:5000/api/v1/passengers?email=" + email
		} else if !isPassenger {
			url = "http://localhost:5001/api/v1/drivers?email=" + email
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
		}

		// Convert salt from hex string to byte array
		decodedSalt, err := hex.DecodeString(salt)
		if err != nil {
			panic(err)
		}
		// Type assert password into string
		password := authenticationInfo["password"].(string)
		saltedPassword := append([]byte(password), decodedSalt...)

		hashedInput := sha256.Sum256(saltedPassword)
		if fmt.Sprintf("%x", hashedInput) != passHash {
			// return HTTP error here
			return
		}
		token := genJWT(id, email, true)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{"token": token, "isPassenger": isPassenger}
		// Encode map to json string
		jsonResp, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		w.Write(jsonResp)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/login", loginHandler)
	fmt.Println("Authentication Microservice")
	fmt.Println("Listening at port 5003")
	log.Fatal(http.ListenAndServe(":5003", router))
}
