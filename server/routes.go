package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	//URL Routes
	Route{"CreateUrl", "POST", "/urls", createUrl},
	Route{"ReadAllUrls", "GET", "/urls", readUrls},
	Route{"DeleteUrl", "DELETE", "/urls/{id}", deleteUrl},
	Route{"GetUrl", "GET", "/urls/{id}", showUrl},
	Route{"UpdateUrl", "PUT", "/urls/{id}", updateUrl},
	//Server Routes
	Route{"CreateServer", "POST", "/servers", createServer},
	Route{"ReadAllServers", "GET", "/servers", readServers},
	Route{"DeleteServer", "DELETE", "/servers/{id}", deleteServer},
	Route{"GetServer", "GET", "/servers/{id}", showServer},
	Route{"UpdateServer", "PUT", "/servers/{id}", updateServer},
	//Stat Routes
	Route{"GetStats", "POST", "/stats/{id}", showStats},
	//User Routes
	Route{"CreateUser", "POST", "/users", createUser},
	Route{"ReadAllUsers", "GET", "/users", readUsers},
	Route{"DeleteUser", "DELETE", "/users/{id}", deleteUser},
	Route{"GetUser", "GET", "/users/{id}", showUser},
	Route{"UpdateUser", "PUT", "/users/{id}", updateUser},
}
