package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "github.com/eminetto/clean-architecture-go/migrations"
)

func main() {
	setupAndRunUsersServer()

	//http.ListenAndServe(":3000", r)
}

func setupAndRunUsersServer() {

	router := mux.NewRouter()

	SetupDatabase()

	router.HandleFunc("/api/users", GetUsers).Methods("GET")
	router.HandleFunc("/api/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/api/users", AddUser).Methods("POST")
	router.HandleFunc("/api/users/{id}", UpdateUser).Methods("PUT")
	router.HandleFunc("/api/users/{id}", DeleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3000", router))
}

var users []User

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Role     string `json:"role"`
}

func SetupDatabase() []User {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://kevs:password@mongorestapi:27017/miapp?authSource=admin"))
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	usersCollection := client.Database("classroomdb").Collection("projectLab")

	collection, err := usersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	} else {
		for collection.Next(ctx) {
			var result User
			err := collection.Decode(&result)
			if err != nil {
				fmt.Println("cursor.Next() error:", err)
				os.Exit(1)
			} else {
				users = append(users, result)
			}
		}
	}
	return users
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	users = append(users, user)
	json.NewEncoder(w).Encode(user)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, user := range users {
		id := (params["id"])
		if user.ID == id {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	json.NewEncoder(w).Encode(&User{})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, value := range users {
		id := params["id"]
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, user := range users {
		id := params["id"]
		if user.ID == id {
			users = append(users[:index], users[index+1:]...)
		}
	}

}
