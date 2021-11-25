package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)


type trip struct{
	ID int
	PickUp string
	DropOff string
	PassengerID int
	DriverID int
	Start time.Time
	End time.Time
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


// Request - POST
// View - GET
// Start - PATCH?
// End - PATCH?
func tripHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method=="GET" {
		kv := r.URL.Query()
		passenger_id := kv["passenger_id"]
		var query string
		if passenger_id!=nil{
			query = fmt.Sprintf("select * from trip where passenger_id="+passenger_id[0])
		}
		results, err := database.Query(query)
		if err != nil {
			panic(err.Error())
		}
		var trips []trip
		for results.Next() {
			// map this type  to the record in the table
			var trip trip
			err = results.Scan(&trip.ID, &trip.PickUp, &trip.DriverID, &trip.PassengerID, &trip.DriverID, &trip.Start, &trip.End)
	
			if err != nil {
				panic(err.Error())
			}
			trips = append(trips,trip)
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
	
			query:=fmt.Sprintf("INSERT INTO trip (pickup,dropoff,passenger_id) VALUES ('%s','%s',%d)",trip.PickUp,trip.DropOff,id)
	
			_, err = database.Query(query)
			if err != nil {
				panic(err.Error())
			}
			w.Write([]byte("200 - Trip created"))
		}
		if r.Method == "PATCH" {
			var trip trip
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("422 - Please supply course information in JSON format"))
				return
			}
			// authenticate user
			id, isPassenger, authenticated := authentication(r)
			// Only drivers can update
			if !authenticated || isPassenger {
				// Return HTTP error here
				return
			}
	
			json.Unmarshal(reqBody, &trip)
			// GET trip
			query := fmt.Sprintf("select * from trip where id="+string(trip.ID))
			results, err := database.Query(query)
			if err != nil {
				panic(err.Error())
			}
			results.Next()
			err = results.Scan(&trip.ID, &trip.PickUp, &trip.DriverID, &trip.PassengerID, &trip.DriverID, &trip.Start, &trip.End)

			if trip.DriverID != id {
				// Return unauthorized status code
				return
			}



			query = ""
			// Only start or end can be updated
			if trip.Start.String() != ""{
				query = fmt.Sprintf("UPDATE trip SET start='%s'",trip.Start.String())
			} else if trip.End.String() != "" {
				query = fmt.Sprintf("UPDATE trip SET end='%s'",trip.End.String())
			}
	
			_, err = database.Query(query)
			if err != nil {
				panic(err.Error())
			}
			w.Write([]byte("200 - Trip updated"))
		}
	}
}


func main(){
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/ooper")
	//  handle error
	if err != nil {
		panic(err.Error())
	}
	database = db

	// defer the close till after the main function has finished  executing
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/passengers", tripHandler).Methods(http.MethodPatch, http.MethodPost, http.MethodOptions, http.MethodGet)
	router.Use(mux.CORSMethodMiddleware(router))
	fmt.Println("Listening at port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}