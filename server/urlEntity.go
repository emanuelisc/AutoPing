package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
	// "fmt"
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

	ua := r.Header.Get("Content-Type")

	if ua != "application/json" {
		responseCode(w, http.StatusUnsupportedMediaType)
		return
	}

	// log.Printf(data)
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

	responseJSONCode(w, url, http.StatusCreated)
}

func readUrls(w http.ResponseWriter, r *http.Request) {
	result := []Url{}
	if err := urls.Find(nil).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
	} else {
		responseJSONCode(w, result, http.StatusOK)
	}
}

func deleteUrl(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseCode(w, http.StatusNotFound)
		return
	}

	if err := urls.RemoveId(bson.ObjectIdHex(params["id"])); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
		return
	}
	responseCode(w, http.StatusNoContent)
}

func showUrl(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseCode(w, http.StatusNotFound)
		return
	}

	result := Url{}
	err := urls.Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).One(&result)
	if err != nil {
		responseError(w, "Invalid Url ID", http.StatusNotFound)
		return
	}
	responseJSONCode(w, result, http.StatusOK)
}

func updateUrl(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseCode(w, http.StatusNotFound)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	ua := r.Header.Get("Content-Type")

	if ua != "application/json" {
		responseCode(w, http.StatusBadRequest)
		return
	}

	url := &Url{}
	err = json.Unmarshal(data, url)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := urls.UpdateId(bson.ObjectIdHex(params["id"]), url); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
		return
	}

	responseJSONCode(w, url, http.StatusOK)
}
