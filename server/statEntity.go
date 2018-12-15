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

func avarageResponseTime(a []float64) float64 {
	sum := float64(0)
	length := float64(0)
	for _, value := range a {
		sum = sum + value
		length = length + 1
	}
	return sum / length
}

func getAvarageResults(data float64) string {
	if data < 0.001 {
		return "Good"
	} else if data > 0.001 && data < 0.01 {
		return "Avarage"
	} else if data > 0.01 && data < 1.0 {
		return "Bad"
	} else {
		return "NoData"
	}
}
