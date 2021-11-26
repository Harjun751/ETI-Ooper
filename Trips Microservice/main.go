package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
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

func authentication(r *http.Request) (int, bool, bool) {
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

func startTripHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == "POST" {
		// authenticate user
		id, isPassenger, authenticated := authentication(r)
		// Only drivers can update
		if !authenticated || isPassenger {
			// Return HTTP error here
			return
		}

		params := mux.Vars(r)
		tripID := params["ID"]
		var driver_id int

		// GET trip
		query := fmt.Sprintf("select driver_id from trip where id=" + tripID)
		results, err := database.Query(query)
		if err != nil {
			panic(err.Error())
		}
		results.Next()
		err = results.Scan(&driver_id)
		if err != nil {
			panic(err)
		}
		if driver_id != id {
			// Return unauthorized status code
			return
		}
		_, err = database.Query("UPDATE trip SET start=NOW() where id=" + tripID)
		if err != nil {
			panic(err.Error())
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
		id, isPassenger, authenticated := authentication(r)
		// Only drivers can update
		if !authenticated || isPassenger {
			// Return HTTP error here
			return
		}

		params := mux.Vars(r)
		tripID := params["ID"]
		var driver_id int

		// GET trip
		query := fmt.Sprintf("select driver_id from trip where id=" + tripID)
		results, err := database.Query(query)
		if err != nil {
			panic(err.Error())
		}
		results.Next()
		err = results.Scan(&driver_id)
		if err != nil {
			panic(err)
		}
		if driver_id != id {
			// Return unauthorized status code
			return
		}
		_, err = database.Query("UPDATE trip SET end=NOW() where id=" + tripID)
		if err != nil {
			panic(err.Error())
		}

		// Update availability of driver
		url := fmt.Sprintf("http://localhost:5001/api/v1/drivers/%d/availability", driver_id)
		newReqBody, err := json.Marshal(map[string]interface{}{"availability": true})
		if err != nil {
			panic(err.Error())
		}
		request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(newReqBody))
		if err != nil {
			// Fail to set availability
			// TODO: Log
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
		id, isPassenger, authenticated := authentication(r)
		if !authenticated {
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
			panic(err.Error())
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

	if r.Method == "GET" {
		id, _, authenticated := authentication(r)
		if !authenticated {
			return
		}
		query := fmt.Sprintf("select * from trip where passenger_id=%d", id)
		results, err := database.Query(query)
		if err != nil {
			panic(err.Error())
		}
		var trips []trip
		for results.Next() {
			// map this type  to the record in the table
			var trip trip
			err = results.Scan(&trip.ID, &trip.PickUp, &trip.DropOff, &trip.PassengerID, &trip.DriverID, &trip.Requested, &trip.Start, &trip.End)

			if err != nil {
				panic(err.Error())
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
			// authenticate user
			id, isPassenger, authenticated := authentication(r)
			if !authenticated || !isPassenger {
				// Return HTTP error here
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
			resp, err := http.Get("http://localhost:5001/api/v1/drivers/available")
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
				panic(err.Error())
			}

			// Update availability of driver
			url := fmt.Sprintf("http://localhost:5001/api/v1/drivers/%d/availability", driver_id)
			newReqBody, err := json.Marshal(map[string]interface{}{"availability": false})
			if err != nil {
				panic(err.Error())
			}
			request, err := http.NewRequest("PATCH", url, bytes.NewBuffer(newReqBody))
			if err != nil {
				// Fail to set availability
				// TODO: Log
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
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ooper")
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