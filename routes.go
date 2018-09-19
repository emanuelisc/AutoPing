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
	Route{"ReadUrls", "GET", "/urls", readUrls},
	Route{"DeleteUrl", "DELETE", "/urls/{id}", deleteUrl},
	Route{"GetUrl", "GET", "/urls/{id}", showUrl},
	//Server Routes
	Route{"CreateServer", "POST", "/servers", createServer},
	Route{"ReadServers", "GET", "/servers", readServers},
	Route{"DeleteServer", "DELETE", "/servers/{id}", deleteServer},
	Route{"GetServer", "GET", "/servers/{id}", showServer},
}
