package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	"strings"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type User struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Servers     string        `json:"servers" bson:"servers"`
	Username		string        `json:"username" bson:"username"`
	Password		string        `json:"password" bson:"password"`
	CreatedAt   time.Time     `json:"createdAt" bson:"created_at"`
}

func createUser(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	ua := r.Header.Get("Content-Type")
	if !strings.Contains(ua, "application/json") {
		responseCode(w, http.StatusUnsupportedMediaType)
		return
	}
	user := &User{}
	err = json.Unmarshal(data, user)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return	
	}
	result := User{}
	err = users.Find(bson.M{"username": user.Username}).One(&result)
	log.Printf("1")
	if err == nil {
		responseError(w, "User already exists", http.StatusConflict)
		log.Printf("2")
		return
	}
	log.Printf("4")
	user.CreatedAt = time.Now().UTC()
	password := user.Password
	
	user.Password, err = HashPassword(password)
	if err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if err := users.Insert(user); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseCode(w, http.StatusCreated)
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

	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseCode(w, http.StatusNotFound)
		return
	}

	if err := users.RemoveId(bson.ObjectIdHex(params["id"])); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
		return
	}
	responseCode(w, http.StatusNoContent)
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