package main

import (
	// "encoding/json"
	// "io/ioutil"
	"net/http"
	"time"
	// "fmt"
	// "github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

type Ping struct {
	Id          bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name        string        `json:"name" bson:"name"`
	Ip          string        `json:"ip" bson:"ip"`
	Description string        `json:"description" bson:"description"`
	Servers     string        `json:"servers" bson:"servers"`
	User        string        `json:"user" bson:"user"`
	CreatedAt   time.Time     `json:"createdAt" bson:"created_at"`
}

func pingReq(w http.ResponseWriter, r *http.Request) {
	// result := []Ping{}
	// if err := urls.Find(nil).Sort("-created_at").All(&result); err != nil {
	// 	responseError(w, err.Error(), http.StatusNotFound)
	// } else {
	// 	responseJSONCode(w, result, http.StatusOK)
	// }
}
