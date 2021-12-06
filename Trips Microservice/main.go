package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type trip struct {
	ID          int
	PickUp      string
	DropOff     string
	PassengerID int
	DriverID    int
	Requested   mysql.NullTime
	Start       mysql.NullTime
	End         mysql.NullTime
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

func startTripHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == "POST" {
		// authenticate user
		id, isPassenger, errorStatusCode, errorText := getAuthDetails(r.Header.Get("Authorization"))
		if errorStatusCode != 0 {
			w.WriteHeader(errorStatusCode)
			w.Write([]byte(errorText))
			return
		}
		// Only drivers can update
		if isPassenger {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Not authorized for action"))
			return
		}

		params := mux.Vars(r)
		tripID := params["ID"]
		var driver_id int

		// GET trip
		query := fmt.Sprintf("select driver_id from trip where id=" + tripID)
		results, err := database.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
			return
		}
		results.Next()
		err = results.Scan(&driver_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}
		if driver_id != id {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized to perform action"))
			return
		}
		_, err = database.Query("UPDATE trip SET start=NOW() where id=" + tripID)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
			return
		}
		w.Write([]byte("200 - Trip updated"))
	}

}
func endTripHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == "POST" {
		// authenticate user
		id, isPassenger, errorStatusCode, errorText := getAuthDetails(r.Header.Get("Authorization"))
		if errorStatusCode != 0 {
			w.WriteHeader(errorStatusCode)
			w.Write([]byte(errorText))
			return
		}
		// Only drivers can update
		if isPassenger {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized to perform action"))
			return
		}

		params := mux.Vars(r)
		tripID := params["ID"]
		var driver_id int

		// GET trip
		query := fmt.Sprintf("select driver_id from trip where id=" + tripID)
		results, err := database.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
			return
		}
		results.Next()
		err = results.Scan(&driver_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}
		if driver_id != id {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized to perform action"))
			return
		}
		_, err = database.Query("UPDATE trip SET end=NOW() where id=" + tripID)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
			return
		}

		// Update availability of driver
		url := fmt.Sprintf(os.Getenv("DRIVER_MS_HOST")+"/api/v1/drivers/%d/availability", driver_id)
		newReqBody, err := json.Marshal(map[string]interface{}{"availability": true})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Error"))
			return
		}
		request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(newReqBody))
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Driver Endpoint Unavailable"))
			return
		}
		request.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		client.Do(request)

		w.Write([]byte("200 - Trip updated"))
	}

}

func currentTripHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == "GET" {
		id, isPassenger, errorStatusCode, errorText := getAuthDetails(r.Header.Get("Authorization"))
		if errorStatusCode != 0 {
			w.WriteHeader(errorStatusCode)
			w.Write([]byte(errorText))
			return
		}
		var query string
		if isPassenger {
			query = fmt.Sprintf("select * from trip where passenger_id=%d order by requested desc limit 1", id)
		} else {
			query = fmt.Sprintf("select * from trip where driver_id=%d and end is null order by requested desc limit 1;", id)
		}
		results := database.QueryRow(query)
		var trip trip
		err := results.Scan(&trip.ID, &trip.PickUp, &trip.DropOff, &trip.PassengerID, &trip.DriverID, &trip.Requested, &trip.Start, &trip.End)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - No current trip"))
			return
		}
		json.NewEncoder(w).Encode(trip)
	}
}

func tripHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}
	// authenticate user
	// Both GET and POST require authentication so do it here
	id, isPassenger, errorStatusCode, errorText := getAuthDetails(r.Header.Get("Authorization"))
	if errorStatusCode != 0 {
		w.WriteHeader(errorStatusCode)
		w.Write([]byte(errorText))
		return
	}

	if r.Method == "GET" {
		query := fmt.Sprintf("select * from trip where passenger_id=%d", id)
		results, err := database.Query(query)
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 - Database Unavailable"))
			return
		}
		var trips []trip
		for results.Next() {
			// map this type  to the record in the table
			var trip trip
			err = results.Scan(&trip.ID, &trip.PickUp, &trip.DropOff, &trip.PassengerID, &trip.DriverID, &trip.Requested, &trip.Start, &trip.End)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Error"))
				return
			}
			trips = append(trips, trip)
		}

		json.NewEncoder(w).Encode(trips)
	}

	if r.Header.Get("Content-Type") == "application/json" {
		if r.Method == "POST" {
			var trip trip
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}
			if !isPassenger {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 - Unauthorized to perform action"))
				return
			}

			json.Unmarshal(reqBody, &trip)
			if trip.DropOff == "" || trip.PickUp == "" {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}

			var driver_id int
			var first_name string
			var last_name string
			var license_number string
			// Get available driver
			resp, err := http.Get(os.Getenv("DRIVER_MS_HOST") + "/api/v1/drivers/available")
			if err == nil {
				defer resp.Body.Close()
				if body, err := ioutil.ReadAll(resp.Body); err == nil {
					var result map[string]interface{}
					json.Unmarshal(body, &result)
					driver_id = int(result["ID"].(float64))
					first_name = result["FirstName"].(string)
					last_name = result["LastName"].(string)
					license_number = result["LicenseNumber"].(string)
				}
			}

			query := fmt.Sprintf("INSERT INTO trip (pickup,dropoff,passenger_id,driver_id,requested) VALUES ('%s','%s',%d,%d,NOW())", trip.PickUp, trip.DropOff, id, driver_id)

			_, err = database.Query(query)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("503 - Database Unavailable"))
				return
			}

			// Update availability of driver
			url := fmt.Sprintf(os.Getenv("DRIVER_MS_HOST")+"/api/v1/drivers/%d/availability", driver_id)
			newReqBody, err := json.Marshal(map[string]interface{}{"availability": false})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("500 - Internal Error"))
				return
			}
			request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(newReqBody))
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("503 - Driver Endpoint Unavailable"))
				return
			}
			request.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			client.Do(request)

			// return driver
			json.NewEncoder(w).Encode(map[string]interface{}{"FirstName": first_name, "LastName": last_name, "LicenseNumber": license_number})
		}
	}
}

func main() {
	if os.Getenv("ENVIRONMENT") != "production" {
		os.Setenv("MYSQL_HOST", "localhost:3306")
		os.Setenv("DATABASE_NAME", "ooper")
		os.Setenv("AUTH_MS_HOST", "http://localhost:5003")
		os.Setenv("DRIVER_MS_HOST", "http://localhost:5001")
		fmt.Println("Using localhost:3306 as database host and ooper as database name")
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

	router.HandleFunc("/api/v1/trips", tripHandler).Methods(http.MethodPost, http.MethodOptions, http.MethodGet)
	router.HandleFunc("/api/v1/trip/{ID}/start", startTripHandler).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/trip/{ID}/end", endTripHandler).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/api/v1/current-trip", currentTripHandler).Methods(http.MethodOptions, http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	fmt.Println("Trips Microservice")
	fmt.Println("Listening at port 5004")
	log.Fatal(http.ListenAndServe(":5004", router))
}
