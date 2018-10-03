package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Ip          string        `json:"ip" bson:"ip"`
	Description string        `json:"description" bson:"description"`
	Servers     string        `json:"servers" bson:"servers"`
	User        string        `json:"user" bson:"user"`
	CreatedAt   time.Time     `json:"createdAt" bson:"created_at"`
}

func createUser(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &User{}
	err = json.Unmarshal(data, user)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.CreatedAt = time.Now().UTC()

	if err := users.Insert(user); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, user)
}

func readUsers(w http.ResponseWriter, r *http.Request) {
	result := []User{}
	if err := users.Find(nil).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseJSON(w, result)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	if err := users.RemoveId(bson.ObjectIdHex(params["id"])); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseJSON(w, http.StatusOK)
}

func showUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result := User{}
	err := users.Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).One(&result)
	if err != nil {
		responseError(w, "Invalid User ID", http.StatusBadRequest)
		return
	}
	responseJSON(w, result)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &User{}
	err = json.Unmarshal(data, user)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := users.UpdateId(bson.ObjectIdHex(params["id"]), user); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, user)
}
