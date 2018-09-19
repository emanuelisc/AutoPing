package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type Url struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Ip          string        `json:"ip" bson:"ip"`
	Description string        `json:"description" bson:"description"`
	Servers     string        `json:"servers" bson:"servers"`
	User        string        `json:"user" bson:"user"`
	CreatedAt   time.Time     `json:"createdAt" bson:"created_at"`
}

func createUrl(w http.ResponseWriter, r *http.Request) {

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := &Url{}
	err = json.Unmarshal(data, url)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	url.CreatedAt = time.Now().UTC()

	if err := urls.Insert(url); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSON(w, url)
}

func readUrls(w http.ResponseWriter, r *http.Request) {
	// Get posts collection
	result := []Url{}
	if err := urls.Find(nil).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
	} else {
		responseJSON(w, result)
	}
}

func deleteUrl(w http.ResponseWriter, r *http.Request) {
	// var url Url
	params := mux.Vars(r)
	// result := Url{}
	// err := urls.Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).One(&result)
	// if err != nil {
	// 	responseError(w, "Invalid Url ID", http.StatusBadRequest)
	// 	return
	// }

	if err := urls.RemoveId(bson.ObjectIdHex(params["id"])); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseJSON(w, http.StatusOK)
}

func showUrl(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result := Url{}
	err := urls.Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).One(&result)
	if err != nil {
		responseError(w, "Invalid Url ID", http.StatusBadRequest)
		return
	}
	responseJSON(w, result)
}
