package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// swagger:model
type User struct {

	// The ID
	// min: 0
	// max: 30
	ID int `json:"id"`

	// The name for this user
	// required: true
	// min lenght: 1
	Name string `json:"name"`
}

var users []User
var idCounter int

func main() {
	runServer()

}

func runServer() {
	router := mux.NewRouter()

	database()

	router.HandleFunc("/api/users", getUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/api/users", addUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", deleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))

}

func database() {
	users = append(users, User{ID: 1, Name: "Kevin"})
	users = append(users, User{ID: 2, Name: "Oscar"})
	idCounter = 2
}
func addUser(w http.ResponseWriter, r *http.Request) {

	// swagger:operation POST /users postUser
	//
	// Include documentation
	//
	// ---
	// produces:
	// - application/json
	// parameters:
	//   - name: Body
	//     in: body
	//     schema:
	//       "$ref": "#/definitions/User"
	// responses:
	//   '200':
	//     description: user response

	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	idCounter++
	user.ID = idCounter
	users = append(users, user)

	json.NewEncoder(w).Encode(user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /users getUsers
	//
	// Insert documentation
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: user response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/User"

	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	// swagger:operation GET /users getUsers
	//
	// Insert documentation
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: user response
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/User"

	params := mux.Vars(r)
	for _, user := range users {
		id, _ := strconv.Atoi(params["id"])

		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

func updateUser(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	for index, value := range users {
		id, _ := strconv.Atoi(params["id"])
		if value.ID == id {
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)
			user.ID = users[index].ID
			users[index].Name = user.Name
			json.NewEncoder(w).Encode(&User{})
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, user := range users {

		id, _ := strconv.Atoi(params["id"])
		if user.ID == id {
			users = append(users[:index], users[index+1:]...)
		}
	}

}
