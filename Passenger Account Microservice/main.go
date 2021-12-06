package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

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
	Salt         string
}

var database *sql.DB

func getAuthDetails(header string) (id int, isPassenger bool, errorStatusCode int, errorText string) {
	errorStatusCode = 0
	errorText = ""
	newReqBody, err := json.Marshal(map[string]interface{}{"authorization": header})
	if err != nil {
		errorStatusCode = http.StatusInternalServerError
		errorText = "500 - Internal Error"
		return
	}
	// POST to authentication microservice with details
	resp, err := http.Post(os.Getenv("AUTH_MS_HOST")+"/api/v1/authorize", "application/json", bytes.NewBuffer(newReqBody))
	if err == nil {
		if resp.StatusCode != 200 {
			errorStatusCode = http.StatusUnprocessableEntity
			errorText = "401 - Access Token Incorrect"
			return
		}
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			var result map[string]interface{}
			json.Unmarshal(body, &result)
			id = int(result["ID"].(float64))
			isPassenger = result["isPassenger"].(bool)
		}
	} else if err != nil {
		errorStatusCode = http.StatusServiceUnavailable
		errorText = "503 - Authentication unavailable"
		return
	}
	return
}

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
func passengersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == "GET" {
		kv := r.URL.Query()
		id := kv["id"]
		email := kv["email"]
		var query string
		if id != nil {
			query = fmt.Sprintf("select * from passenger where ID=%s", id[0])
		} else if email != nil {
			query = fmt.Sprintf("select * from passenger where email='%s'", email[0])
		}
		results, err := database.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
			log.Print(err)
			return
		}
		results.Next()
		var passenger passenger
		err = results.Scan(&passenger.ID, &passenger.FirstName, &passenger.LastName, &passenger.MobileNumber, &passenger.Email, &passenger.Password, &passenger.Salt)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
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
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("503 - Database Unavailable"))
				log.Print(err)
				return
			}
			w.Write([]byte("200 - Account created"))
		}

		if r.Method == "PATCH" {
			var newPassenger passenger
			reqBody, err := ioutil.ReadAll(r.Body)
			// authenticate user
			id, isPassenger, errorStatusCode, errorText := getAuthDetails(r.Header.Get("Authorization"))
			if errorStatusCode != 0 {
				w.WriteHeader(errorStatusCode)
				w.Write([]byte(errorText))
				return
			}
			if !isPassenger {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("401 - Access token incorrect/unauthorized"))
				return
			}

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

			query := fmt.Sprintf("UPDATE passenger SET first_name='%s',last_name='%s',mobile_number=%d,email='%s' WHERE ID=%d;", newPassenger.FirstName, newPassenger.LastName, newPassenger.MobileNumber, newPassenger.Email, id)
			_, err = database.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("503 - Database Unavailable"))
				log.Print(err)
				return
			}
		}
	}
}

func main() {
	if os.Getenv("ENVIRONMENT") != "production" {
		os.Setenv("MYSQL_HOST", "localhost:3306")
		os.Setenv("DATABASE_NAME", "ooper")
		os.Setenv("AUTH_MS_HOST", "http://localhost:5003")
	}
	db, err := sql.Open("mysql", "user:password@tcp("+os.Getenv("MYSQL_HOST")+")/"+os.Getenv("DATABASE_NAME"))

	//  handle error
	if err != nil {
		panic(err.Error())
	}
	database = db

	// defer the close till after the main function has finished  executing
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/passengers", passengersHandler).Methods(http.MethodPatch, http.MethodPost, http.MethodOptions, http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	fmt.Println("Passenger Microservice")
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
