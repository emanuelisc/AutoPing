package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	// "fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
)

type Server struct {
	Id        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name      string        `json:"name" bson:"name"`
	HostName  string        `json:"hostname" bson:"hostname"`
	Ip        string        `json:"ip" bson:"ip"`
	Status    bool          `json:"status" bson:"status"`
	CreatedAt time.Time     `json:"createdAt" bson:"created_at"`
}

func createServer(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in User
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	// 2. Check user role

	// 3. Check Content-Type
	if r.Header.Get("Content-Type") != "application/json" {
		responseCode(w, http.StatusBadRequest)
		return
	}

	// 4. Get Data from Response body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 5. Get an Object from Data
	server := &Server{}
	err = json.Unmarshal(data, server)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 6. Add few other fields to the Object
	server.CreatedAt = time.Now().UTC()

	// 7. Insert an Object to the Database
	if err := servers.Insert(server); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSONCode(w, server, http.StatusCreated)
}

func readServers(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in User
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	// 2. Check User role

	// 3. Get Urls
	result := []Server{}
	if err := servers.Find(nil).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
	} else {
		responseJSONCode(w, result, http.StatusOK)
	}
}

func deleteServer(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in User
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	// 2. Check User role

	// 3. Get Params
	params := mux.Vars(r)

	// 4. Check if Params are valid
	id := bson.ObjectIdHex(params["id"])
	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseJSONCode(w, http.StatusNotFound, http.StatusNotFound)
		return
	}

	// 5. Remove Object from Database
	if err := servers.RemoveId(id); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
		return
	}
	responseJSON(w, http.StatusOK)
}

func showServer(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in User
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	// 2. Check User role

	// 3. Get Params
	params := mux.Vars(r)

	// 4. Check if Params are valid
	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseJSONCode(w, http.StatusNotFound, http.StatusNotFound)
		return
	}

	// 5. Find Object 
	result := Server{}
	err := servers.Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).One(&result)
	if err != nil {
		responseError(w, "Invalid Server ID", http.StatusNotFound)
		return
	}
	responseJSON(w, result)
}

func updateServer(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in user
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	// 2. Check User role

	// 3. Check Content-Type
	ua := r.Header.Get("Content-Type")
	if ua != "application/json" || ua != "application/json; charset=utf-8" {
		responseCode(w, http.StatusBadRequest)
		return
	}

	// 3. Get Params
	params := mux.Vars(r)

	// 4. Check if Params are valid
	valid := bson.IsObjectIdHex(params["id"])
	id := bson.ObjectIdHex(params["id"])
	if valid != true {
		responseCode(w, http.StatusNotFound)
		return
	}

	// 5. Get Data from Responce body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 6. Get Object from Data
	server := &Server{}
	err = json.Unmarshal(data, server)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 7. Return Data
	if err := servers.UpdateId(id, server); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
		return
	}

	responseJSONCode(w, server, http.StatusOK)
}
