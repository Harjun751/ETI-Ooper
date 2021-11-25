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

		var id int
		email := authenticationInfo["email"].(string)
		var salt string
		var passHash string
		if resp, err := http.Get("localhost:5000/api/v1/passengers?email="+email); err == nil {
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				var result map[string]interface{}
				json.Unmarshal(body, &result)
				if result["success"]==true{
					id = int(result["id"].(float64))
					salt = result["salt"].(string)
					passHash = result["passHash"].(string)					
				} else {
					// HTTP error here
					return
				}
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
		// TODO: Change Code so that it allows both drivers and passengers to log in
		// Or create new login endpoint for driver
		token := genJWT(id, email, true)
		w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["token"] = token
		resp["isPassenger"] = "true"
		// Encode map to json string
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		w.Write(jsonResp)
	}
}

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/login", loginHandler)
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}