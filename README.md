# ETI-Ooper

Ooper is a ride sharing platform created for an assignment.

# Microservices & Operations

| Operation      | Description                                                     | Microservice |
| -------------- | --------------------------------------------------------------- | ------------ |
| Create Account | Creates new account object                                      | Accounts     |
| Update Account | Updates account object                                          | Accounts     |
| Login          | Authenticates user                                              | Accounts     |
| Assign Driver  | Given a list of driver objects, selects most appropriate driver | Assigning    |
| Request Trip   | Starts a new Trip for a given customer                          | Trips        |
| View Trip      | Returns a list of trips for a given customer                    | Trips        |
| Start Trip     | Given a trip, updates status of trip                            | Trips        |
| End Trip       | Given a trip, updates status of trip                            | Trips        |

| Microservice | Database | Dependencies                                          |
| ------------ | -------- | ----------------------------------------------------- |
| Accounts     | y        | nil                                                   |
| Assigning    | n        | Accounts - To obtain driver list                      |
| Trips        | y        | Assigning - Automatically invoked after request trip? |
