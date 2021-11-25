package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type driver struct {
	ID           int
	FirstName    string
	LastName     string
	MobileNumber int
	Email        string
	Password     string
	ICNumber 	 string
	LicenseNumber string
	Salt 		 string
}

var database *sql.DB

func authentication(r *http.Request) (int,bool, bool) {
	var secret = []byte("it took the night to believe")
	headerToken := r.Header.Get("Authorization")
	// Decode the jwt and ensure it's readable
	token, err := jwt.Parse(headerToken[7:], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := int(claims["id"].(float64))
		isPassenger := claims["isPassenger"].(bool)
		return id, isPassenger, true
	} else {
		fmt.Println(err)
		return 0, true, false
	}
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
		if id!=nil {
			query = fmt.Sprintf("select * from driver where ID=%s", id[0])
		} else if email != nil{
			query = fmt.Sprintf("select * from driver where email='%s'", email[0])
		}
		results, err := database.Query(query)
		if err != nil {
			panic(err.Error())
		}
		results.Next()
		var driver driver
		err = results.Scan(&driver.ID, &driver.FirstName, &driver.LastName, &driver.MobileNumber, &driver.Email, &driver.ICNumber, &driver.LicenseNumber,&driver.Password, &driver.Salt)

		if err != nil {
			panic(err.Error())
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

			if newDriver.Email == "" || newDriver.FirstName == "" || newDriver.LastName == "" || newDriver.Password == ""|| newDriver.ICNumber == "" || newDriver.LicenseNumber == "" {
				// Return error for incomplete passenger details
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			// get salt and hash of password
			salt, hash := saltNHash(newDriver.Password)

			query := fmt.Sprintf("INSERT INTO driver (first_name,last_name,mobile_number,email,ic_number,license_number,password,salt) VALUES ('%s', '%s', %d, '%s', '%s', '%s','%s','%s')", newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.Email, newDriver.ICNumber,newDriver.LicenseNumber,hash, salt)
			_, err = database.Query(query)
			if err != nil {
				panic(err.Error())
			}
			w.Write([]byte("200 - Account created"))
		}

		if r.Method == "PATCH" {
			var newDriver driver
			reqBody, err := ioutil.ReadAll(r.Body)
			// authenticate user
			id, isPassenger, authenticated := authentication(r)
			if !authenticated || isPassenger {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("401 - Access token incorrect/unauthorized"))
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

			query := fmt.Sprintf("UPDATE driver SET first_name='%s',last_name='%s',mobile_number=%d,email='%s',ic_number='%s',license_nmber='%s' WHERE ID=%d;", newDriver.FirstName, newDriver.LastName, newDriver.MobileNumber, newDriver.Email,newDriver.ICNumber,newDriver.LicenseNumber, id)
			_, err = database.Query(query)
			if err != nil {
				panic(err.Error())
			}
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

	router.HandleFunc("/api/v1/drivers", driversHandler).Methods(http.MethodPatch, http.MethodPost, http.MethodOptions, http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	fmt.Println("Listening at port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))
}
