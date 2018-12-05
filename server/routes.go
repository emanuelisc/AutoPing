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
	Route{"CreateUrl", "POST", "/urls", ValidateMiddleware(createUrl)},
	Route{"ReadAllUrls", "GET", "/urls", ValidateMiddleware(readUrls)},
	Route{"DeleteUrl", "DELETE", "/urls/{id}", ValidateMiddleware(deleteUrl)},
	Route{"GetUrl", "GET", "/urls/{id}", ValidateMiddleware(showUrl)},
	Route{"UpdateUrl", "PUT", "/urls/{id}", ValidateMiddleware(updateUrl)},
	//Server Routes
	Route{"CreateServer", "POST", "/servers", ValidateMiddleware(createServer)},
	Route{"ReadAllServers", "GET", "/servers", readServers},
	Route{"DeleteServer", "DELETE", "/servers/{id}", ValidateMiddleware(deleteServer)},
	Route{"GetServer", "GET", "/servers/{id}", showServer},
	Route{"UpdateServer", "PUT", "/servers/{id}", ValidateMiddleware(updateServer)},
	//Stat Routes
	Route{"GetStats", "POST", "/stats/{id}", showStats},
	//Ping Routes
	Route{"GetPing", "GET", "/ping/{address}", pingReq},
	//User Routes
	Route{"CreateUser", "POST", "/users", createUser},
	Route{"ReadAllUsers", "GET", "/users", readUsers},
	Route{"DeleteUser", "DELETE", "/users/{id}", deleteUser},
	Route{"GetUser", "GET", "/users/{id}", showUser},
	Route{"UpdateUser", "PUT", "/users/{id}", updateUser},
	//Authentication
	Route{"Login", "POST", "/authenticate", CreateTokenEndpoint},
	Route{"Protected", "GET", "/protected", ProtectedEndpoint},
	Route{"Test", "GET", "/test", ValidateMiddleware(TestEndpoint)},
}
