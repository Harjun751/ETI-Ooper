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
	"strconv"

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

func getAuthDetails(header string) (id int, isPassenger bool, errorStatusCode int, errorText string){
	errorStatusCode = 0
	errorText = ""
	newReqBody, err := json.Marshal(map[string]interface{}{"authorization": header})
	if err != nil {
		errorStatusCode = http.StatusInternalServerError
		errorText = "500 - Internal Error"
		return
	}
	// POST to authentication microservice with details
	resp, err := http.Post("http://localhost:5003/api/v1/authorize","application/json",bytes.NewBuffer(newReqBody))
	if err == nil {
		if (resp.StatusCode!=200){
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

func setAvailabilityDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == "PATCH" {
		params := mux.Vars(r)
		id := params["ID"]

		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply course information in JSON format"))
			return
		}

		var bodyData map[string]interface{}
		json.Unmarshal(reqBody, &bodyData)

		availability := bodyData["availability"].(bool)
		ID, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Malformed body"))
			return
		}
		query := fmt.Sprintf("UPDATE driver set available=%t where id=%d", availability, ID)
		_, err = database.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
			return
		}
		w.Write([]byte("200 - Updated"))
	}
}

func getAvailableDriver(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == "GET" {
		results, err := database.Query("select id,first_name,last_name,license_number from driver where available=true limit 1;")
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
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
	}
}

func driversHandler(w http.ResponseWriter, r *http.Request) {
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
			query = fmt.Sprintf("select * from driver where ID=%s", id[0])
		} else if email != nil {
			query = fmt.Sprintf("select * from driver where email='%s'", email[0])
		}
		results, err := database.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
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

			// get salt and hash of password
			salt, hash := saltNHash(newDriver.Password)

			query := fmt.Sprintf("INSERT INTO driver (first_name,last_name,mobile_number,email,ic_number,license_number,password,salt) VALUES ('%s', '%s', %d, '%s', '%s', '%s','%s','%s')", newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.Email, newDriver.ICNumber, newDriver.LicenseNumber, hash, salt)
			_, err = database.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("503 - Database Unavailable"))
				return
			}
			w.Write([]byte("200 - Account created"))
		}

		if r.Method == "PATCH" {
			var newDriver driver
			reqBody, err := ioutil.ReadAll(r.Body)
			// authorize user - obtain jwt details from auth microservice
			id, isPassenger, errorStatusCode, errorText := getAuthDetails(r.Header.Get("Authorization"))
			if (errorStatusCode != 0){
				w.WriteHeader(errorStatusCode)
				w.Write([]byte(errorText))
				return
			}
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

			query := fmt.Sprintf("UPDATE driver SET first_name='%s',last_name='%s',mobile_number=%d,email='%s',license_number='%s' WHERE ID=%d;", newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.Email, newDriver.LicenseNumber, id)
			_, err = database.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("503 - Database Unavailable"))
				return
			}
		}
	}
}

func main() {
	db, err := sql.Open("mysql", "user:password@tcp("+os.Getenv("MYSQL_HOST")+")/ooper")

	//  handle error
	if err != nil {
		panic(err.Error())
	}
	database = db

	// defer the close till after the main function has finished  executing
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/drivers", driversHandler).Methods(http.MethodPatch, http.MethodPost, http.MethodOptions, http.MethodGet)
	router.HandleFunc("/api/v1/drivers/available", getAvailableDriver).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/drivers/{ID}/availability", setAvailabilityDriver).Methods(http.MethodPatch)
	router.Use(mux.CORSMethodMiddleware(router))
	fmt.Println("Driver Microservice")
	fmt.Println("Listening at port 5001")
	fmt.Println(os.Getenv("MYSQL_HOST"))
	fmt.Println("HI")
	log.Fatal(http.ListenAndServe(":5001", router))
}
