package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
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
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	server := &Server{}
	err = json.Unmarshal(data, server)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}
	server.CreatedAt = time.Now().UTC()

	if err := servers.Insert(server); err != nil {
		responseError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseJSONCode(w, server, http.StatusCreated)
}

func readServers(w http.ResponseWriter, r *http.Request) {
	result := []Server{}
	if err := servers.Find(nil).Sort("-created_at").All(&result); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
	} else {
		responseJSONCode(w, result, http.StatusOK)
	}
}

func deleteServer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseJSONCode(w, http.StatusNotFound, http.StatusNotFound)
		return
	}

	if err := servers.RemoveId(bson.ObjectIdHex(params["id"])); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
		return
	}
	responseJSON(w, http.StatusOK)
}

func showServer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseJSONCode(w, http.StatusNotFound, http.StatusNotFound)
		return
	}

	result := Server{}
	err := servers.Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).One(&result)
	if err != nil {
		responseError(w, "Invalid Server ID", http.StatusNotFound)
		return
	}
	responseJSON(w, result)
}

func updateServer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	valid := bson.IsObjectIdHex(params["id"])
	if valid != true {
		responseJSONCode(w, http.StatusNotFound, http.StatusNotFound)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	server := &Server{}
	err = json.Unmarshal(data, server)
	if err != nil {
		responseError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := servers.UpdateId(bson.ObjectIdHex(params["id"]), server); err != nil {
		responseError(w, err.Error(), http.StatusNotFound)
		return
	}

	responseJSONCode(w, server, http.StatusOK)
}
