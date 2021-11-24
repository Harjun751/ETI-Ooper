package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type passenger struct {
	ID           int
	FirstName    string
	LastName     string
	MobileNumber int
	Email        string
	Password     string
}

var database *sql.DB

func saltNHash(password string) (string, string) {
	// create a salt of 16 bytes
	salt, _ := GenerateRandomBytes(16)
	// append salt to password
	passwordBytes := append([]byte(password), salt...)
	// create hash from password
	hash := sha256.Sum256(passwordBytes)
	// Convert hash and salt from hexadecimal to string to be stored
	hashString := fmt.Sprintf("%x", hash)
	saltString := fmt.Sprintf("%x", salt)
	return saltString, hashString
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

func passengerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	params := mux.Vars(r)

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == "GET" {
		id := params["ID"]
		query := fmt.Sprintf("select * from passenger where ID=%s", id)
		results, err := database.Query(query)
		if err != nil {
			panic(err.Error())
		}
		results.Next()
		var passenger passenger
		err = results.Scan(&passenger.ID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNumber, &passenger.Email)

		if err != nil {
			panic(err.Error())
		}

		json.NewEncoder(w).Encode(passenger)
	}

	if r.Header.Get("Content-Type") == "application/json" {
		if r.Method == "POST" {
			var newPassenger passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			json.Unmarshal(reqBody, &newPassenger)

			if newPassenger.Email == "" || newPassenger.FirstName == "" || newPassenger.LastName == "" || newPassenger.Password == "" {
				// Return error for incomplete passenger details
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			// get salt and hash of password
			salt, hash := saltNHash(newPassenger.Password)

			query := fmt.Sprintf("INSERT INTO passenger (first_name,last_name,mobile_number,email,password,salt) VALUES ('%s', '%s', %d, '%s', '%s', '%s')", newPassenger.FirstName, newPassenger.LastName, newPassenger.MobileNumber, newPassenger.Email, hash, salt)
			_, err = database.Query(query)
			if err != nil {
				panic(err.Error())
			}
		}

		if r.Method == "PATCH" {
			var newPassenger passenger
			reqBody, err := ioutil.ReadAll(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			json.Unmarshal(reqBody, &newPassenger)

			if newPassenger.Email == "" || newPassenger.FirstName == "" || newPassenger.LastName == "" {
				// Return error for incomplete passenger details
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			query := fmt.Sprintf("UPDATE passenger SET first_name='%s',last_name='%s',mobile_number=%d,email='%s' WHERE ID=%d;", newPassenger.FirstName, newPassenger.LastName, newPassenger.MobileNumber, newPassenger.Email, newPassenger.ID)
			_, err = database.Query(query)
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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

		query := fmt.Sprintf("select password,salt from passenger where email='%s'", authenticationInfo["email"])
		results := database.QueryRow(query)
		var salt string
		var passHash string
		err = results.Scan(&passHash, &salt)
		if err != nil {
			panic(err.Error())
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
		if fmt.Sprintf("%x", hashedInput) == passHash {
			fmt.Println("Authenticated lol")
		} else {
			fmt.Println("Wrong")
		}
	}
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ooper")

	//  handle error
	if err != nil {
		panic(err.Error())
	}
	database = db

	// defer the close till after the main function has finished  executing
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/passengers", passengerHandler).Methods(http.MethodPatch, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/passengers/{ID}", passengerHandler)
	router.HandleFunc("/api/v1/login", loginHandler)
	router.Use(mux.CORSMethodMiddleware(router))
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
