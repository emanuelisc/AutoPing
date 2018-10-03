package main

import (
	// "encoding/json"
	// "io/ioutil"
	"net/http"
	"time"

	// "github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"github.com/sparrc/go-ping"
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

func pingTo(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	// result := Ping{}
	// err := users.Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).One(&result)
	// if err != nil {
	// 	responseError(w, "Invalid User ID", http.StatusBadRequest)
	// 	return
	// }

	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
					panic(err)
	}
	pinger.Count = 3
	pinger.Run() // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats
	// params["target"]
	responseJSON(w, stats)
}
