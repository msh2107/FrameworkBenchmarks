package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type user struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Hobby     string `json:"hobby"`
}

var users []user

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/image", sendImage).Methods("GET")
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./images/"))))

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/user", newUser).Methods("POST")
	r.HandleFunc("/user/{id}", getUser).Methods("GET")
	r.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")

	r.HandleFunc("/sleep", sleep).Methods("GET")

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func sleep(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	_, err := w.Write([]byte("Вы проспали 5 секунд"))
	if err != nil {
		return
	}
}

func sendImage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/images/bug.jpg", http.StatusMovedPermanently)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func newUser(w http.ResponseWriter, r *http.Request) {
	var jsonUser user
	err := json.NewDecoder(r.Body).Decode(&jsonUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	users = append(users, jsonUser)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonUser)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadRequest)
		return
	}
	for _, user := range users {
		if user.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "No such user", http.StatusBadRequest)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadRequest)
		return
	}
	var jsonUser user
	err = json.NewDecoder(r.Body).Decode(&jsonUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for index, user := range users {
		if user.Id == id {
			users[index] = jsonUser
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(jsonUser)
			return
		}
	}
	http.Error(w, "No such user", http.StatusBadRequest)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Bad id", http.StatusBadRequest)
		return
	}
	for index, user := range users {
		if user.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			users = append(users[:index], users[index+1:]...)
			return
		}
	}
	http.Error(w, "No such user", http.StatusBadRequest)
}
