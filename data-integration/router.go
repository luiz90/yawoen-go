package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

var controller = &Controller{Repository: Repository{}}

// Route defines a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes defines the list of routes of our API
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/v1/companies",
		controller.Index,
	},
	Route{
		"AddCompany",
		"POST",
		"/v1/companies",
		controller.AddCompany,
	},
	Route{
		"UpdateCompany",
		"PUT",
		"/v1/companies",
		controller.UpdateCompany,
	},
	Route{
		"DeleteCompany",
		"DELETE",
		"/v1/companies/{id}",
		controller.DeleteCompany,
	},
	Route{
		"MergCompany",
		"POST",
		"/v1/data",
		controller.MergeCompany,
	},
	Route{
	"SearchCompanies",
	"GET",
	"/v1/companies/search",
	controller.SearchCompanies,
	},
}

//NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
