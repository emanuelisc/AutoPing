package main

import (
	// "encoding/json"
	// "io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type Stat struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Ip          string        `json:"ip" bson:"ip"`
	Description string        `json:"description" bson:"description"`
	Servers     string        `json:"servers" bson:"servers"`
	User        string        `json:"user" bson:"user"`
	CreatedAt   time.Time     `json:"createdAt" bson:"created_at"`
}

func showStats(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	result := Stat{}
	err := urls.Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).One(&result)
	if err != nil {
		responseError(w, "Invalid Stat ID", http.StatusBadRequest)
		return
	}
	responseJSON(w, result)
}
