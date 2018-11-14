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

type Url struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Ip          string        `json:"ip" bson:"ip"`
	Description string        `json:"description" bson:"description"`
	Servers     []Server        `json:"servers" bson:"servers"`
	User        string		 		`json:"user" bson:"user"`
	CreatedAt   time.Time     `json:"createdAt" bson:"created_at"`
}

func createUrl(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in User
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	// 2. Check Content-Type
	ua := req.Header.Get("Content-Type")
	log.Print(ua)
	if ua != "application/json" || ua != "application/json; charset=utf-8" {
		responseCode(w, http.StatusUnsupportedMediaType)
		return
	}

	// 3. Get Data from Response body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 4. Get an Object from Response Data
	url := &Url{}
	err = json.Unmarshal(data, url)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 5. Add few other fields to the Object 
	url.CreatedAt = time.Now().UTC()
	url.User = user.Username

	// 6. Insert an Object to the Database
	if err := urls.Insert(url); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseJSONCode(w, url, http.StatusCreated)
}

func readUrls(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in User
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)
	
	// 2. Get Urls
	result := []Url{}
	if err := urls.Find(bson.M{"user": user.Username}).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
	} else {
		responseJSONCode(w, result, http.StatusOK)
	}
}

func deleteUrl(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in User
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)
	
	// 2. Get Params
	params := mux.Vars(r)

	// 3. Check if Params are valid
	id := bson.ObjectIdHex(params["id"])
	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseError(w, "Invalid Url ID", http.StatusNotFound)
		return
	}

	// 4. Find Object 
	result := Url{}
	err := urls.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		responseError(w, "Invalid Url ID", http.StatusNotFound)
		return
	}

	// 5. Check if User can delete Resource
	if result.User != user.Username {
		responseError(w, "You don't have permision!", http.StatusUnauthorized)
		return
	}

	// 6. Remove Object from Database
	if err = urls.RemoveId(id); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
		return
	}
	responseError(w, "URL not found", http.StatusNoContent)
}

func showUrl(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in User
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	// 2. Get Params
	params := mux.Vars(r)

	// 3. Check if Params are valid
	id := bson.ObjectIdHex(params["id"])
	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseError(w, "Invalid Url ID", http.StatusNotFound)
		return
	}

	// 4. Find Object 
	result := Url{}
	err := urls.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		responseError(w, "Invalid Url ID", http.StatusNotFound)
		return
	}

	// 5. Return results
	if user.Username != result.User{
		responseCode(w, http.StatusUnauthorized)
		return
	}
	responseJSONCode(w, result, http.StatusOK)
}

func updateUrl(w http.ResponseWriter, r *http.Request) {

	// 1. Get logged in user
	decoded := context.Get(r, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)

	// 2. Check Content-Type
	ua := r.Header.Get("Content-Type")
	if ua != "application/json" {
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
	url := &Url{}
	err = json.Unmarshal(data, url)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 7. Find Object 
	result := Url{}
	err = urls.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		responseError(w, "Invalid Url ID", http.StatusNotFound)
		return
	}

	// 8. Check if User can update Resource
	if result.User != user.Username {
		responseError(w, "You don't have permision!", http.StatusUnauthorized)
		return
	}

	// 9. Add few other fields to the object
	url.User = user.Username

	// 10. Update Object
	if err := urls.UpdateId(id, url); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
		return
	}

	responseJSONCode(w, url, http.StatusOK)
}
