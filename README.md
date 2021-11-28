# ETI-Ooper

Ooper is a ride sharing platform created for an assignment. It utilizes a microservice architecture coded in Go for the backend of the application, and uses VueJS as the frontend framework.

# Microservices & Operations

|Microservice | Endpoint | Methods | Description|
|-|-|-|-|
|Authentication| /api/v1/login| POST| Authenticates user using hashing and salting, and returns a JWT|
|Driver| /api/v1/drivers| PATCH, POST, GET | Creates a new driver object, Gets a driver, or updates a driver|
|Driver| /api/v1/drivers/available | GET |Obtains an available driver (for a trip)|
|Driver| /api/v1/drivers/{ID}/availability | PATCH | Updates specifically the availability of a driver |
|Passenger| /api/v1/passengers| PATCH, POST, GET | Creates a new passenger object, Gets a passenger, or updates a passenger|
|Trip| /api/v1/trips | POST, GET | Obtains all trips for a passenger, or creates a new trip with an assigned driver|
|Trip| /api/v1/current-trip | GET | Obtains the current trip for a driver (or passenger)|
|Trip| /api/v1/trip/{ID}/start | POST | Sets the "start" attribute for a trip to the server's current time|
|Trip| /api/v1/trip/{ID}/end | POST | Sets the "end" attribute for a trip to the server's current time|

*Most endpoints also accept the OPTIONS method for CORS requests*

# Design Considerations of Microservice
The design of the microservices of the application had some key characteristics of a microservice in mind during development.
* Loosely Coupled
   * Designed to have as little calls to each other as possible
   * However, some do still make calls. These cases are minimized, and only done so when required
   * e.g. Authentication sends a GET request to passenger/driver to get details to verify their identity
* Organized around business capabilities
    * Each microservice represents one aspect of the solution
    * Drivers, Passenger, Trips each relate to a single table in a database

## Authentication
A special consideration was the authentication function. It does not have it's own database, but is it's own microservice because it is more resource-intensive than the usual CRUD operations that the other Microservices carry out. This is due to the process in which the authentication function authenticates users - adding a salt to the password and passing it through a hashing function. Then, the authentication function creates a JWT based on the particulars, and sends the JWT as a response.

Giving the login function the ability to scale independently from other functions would hence be helpful. The load on the passengers and drivers microservice would be light in comparison, as it would only be GETting the details of these users and passing it to the authentication microservice.
