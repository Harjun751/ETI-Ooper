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

type driver struct {
	ID            int
	FirstName     string
	LastName      string
	MobileNumber  int
	Email         string
	Password      string
	ICNumber      string
	LicenseNumber string
	Salt          string
	Available     bool
}

var database *sql.DB

//General function to get authorization details from auth microservice
// Returns error code and string if there's an issue, else returns details
func getAuthDetails(jwt string) (id int, isPassenger bool, errorStatusCode int, errorText string) {
	errorStatusCode = 0
	errorText = ""
	newReqBody, err := json.Marshal(map[string]interface{}{"authorization": jwt})
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

// Makes random bytes of length n
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func driversHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == "GET" {
		kv := r.URL.Query()
		// If "available=true" is in query string,
		// Only returns 1 available driver
		if kv.Get("available") == "true" {
			// Obtain available driver only
			results, err := database.Query("select id,first_name,last_name,license_number from driver where available=true limit 1;")
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("503 - Database Unavailable"))
				log.Print(err)
				return
			}
			res := results.Next()
			if !res {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 - No data"))
				return
			}
			var id int
			var firstName string
			var lastName string
			var licenseNumber string
			err = results.Scan(&id, &firstName, &lastName, &licenseNumber)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Error"))
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{"ID": id, "FirstName": firstName, "LastName": lastName, "LicenseNumber": licenseNumber})
			return
		}
		// Else, obtains a specific user given their ID/email
		id := kv["id"]
		email := kv["email"]
		var results *sql.Rows
		var err error
		if id != nil {
			results, err = database.Query("select * from driver where ID=?", id[0])
		} else if email != nil {
			results, err = database.Query("select * from driver where email=?", email[0])
		}
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
			log.Print(err)
			return
		}
		results.Next()
		var driver driver
		err = results.Scan(&driver.ID, &driver.FirstName, &driver.LastName, &driver.MobileNumber, &driver.Email, &driver.ICNumber, &driver.LicenseNumber, &driver.Password, &driver.Salt, &driver.Available)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}

		json.NewEncoder(w).Encode(driver)
	}

	if r.Header.Get("Content-Type") == "application/json" {
		if r.Method == "POST" {
			var newDriver driver
			reqBody, err := ioutil.ReadAll(r.Body)

			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			json.Unmarshal(reqBody, &newDriver)

			if newDriver.Email == "" || newDriver.FirstName == "" || newDriver.LastName == "" || newDriver.Password == "" || newDriver.ICNumber == "" || newDriver.LicenseNumber == "" {
				// Return error for incomplete passenger details
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			// create salt and hash of password
			salt, hash := saltNHash(newDriver.Password)

			// insert all details
			_, err = database.Query("INSERT INTO driver (first_name,last_name,mobile_number,email,ic_number,license_number,password,salt) VALUES (?, ?, ?, ?, ?, ?,?,?)", newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.Email, newDriver.ICNumber, newDriver.LicenseNumber, hash, salt)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("503 - Database Unavailable"))
				log.Print(err)
				return
			}
			w.Write([]byte("200 - Account created"))
		}

		if r.Method == "PATCH" {
			kv := r.URL.Query()
			// if "availability=true"
			// Only update availability of the driver
			if kv.Get("availability") == "true" {
				reqBody, err := ioutil.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusUnprocessableEntity)
					w.Write([]byte("422 - Please supply course information in JSON format"))
					return
				}

				var bodyData map[string]interface{}
				json.Unmarshal(reqBody, &bodyData)

				// Get availability from body data
				// Only availability is update-able
				availability := bodyData["availability"].(bool)
				ID := bodyData["ID"]
				// Foramt database UPDATE query and send
				_, err = database.Query("UPDATE driver set available=? where id=?", availability, ID)
				if err != nil {
					w.WriteHeader(http.StatusServiceUnavailable)
					w.Write([]byte("503 - Database Unavailable"))
					log.Print(err)
					return
				}
				w.Write([]byte("200 - Updated"))
				return
			}
			// Else, update driver's personal details
			var newDriver driver
			reqBody, err := ioutil.ReadAll(r.Body)
			// Obtain jwt cookie from  response
			jwt, err := r.Cookie("jwt")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 - No authorized cookie"))
				return
			}
			// authorize user - obtain jwt details from auth microservice
			id, isPassenger, errorStatusCode, errorText := getAuthDetails(jwt.Value)
			if errorStatusCode != 0 {
				w.WriteHeader(errorStatusCode)
				w.Write([]byte(errorText))
				return
			}
			// Only non-passengers allowed to edit
			if isPassenger {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("401 - Unauthorized"))
				return
			}

			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			json.Unmarshal(reqBody, &newDriver)

			if newDriver.Email == "" || newDriver.FirstName == "" || newDriver.LastName == "" {
				// Return error for incomplete passenger details
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			// Update driver with all values given
			_, err = database.Query("UPDATE driver SET first_name=?,last_name=?,mobile_number=?,email=?,license_number=? WHERE ID=?;", newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.Email, newDriver.LicenseNumber, id)
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

	router.HandleFunc("/api/v1/drivers", driversHandler).Methods(http.MethodPatch, http.MethodPost, http.MethodOptions, http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	fmt.Println("Driver Microservice")
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}
