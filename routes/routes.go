// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/bounce
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package routes

import (
	"net/http"

	"github.com/apex/log"
)

// Route reflects the patterns and handlers for each supported path in our
// API.
type Route struct {
	Name           string
	Pattern        string
	Description    string
	HandlerFunc    http.HandlerFunc
	AllowedMethods []string
}

// Routes is a collection of defined routes, intended for bulk registration
// and auto-indexing on the landing page for this application
type Routes []Route

// Add appends one or many new routes to the existing Routes collection.
func (rs *Routes) Add(r ...Route) {

	for _, newRoute := range r {
		log.Debugf("Add %s to routes ...", newRoute.Name)
		*rs = append(*rs, newRoute)
	}
}

// RegisterWithServeMux registers each recorded Route with the specified
// HTTP ServeMux.
func (rs *Routes) RegisterWithServeMux(mux *http.ServeMux) {

	// TODO: How would we check for errors registering our route?

	for _, route := range *rs {
		log.Debugf("Register %s with ServeMux ...", route.Name)
		mux.HandleFunc(route.Pattern, route.HandlerFunc)
	}

}

// ListNames provides a list of all recorded route names in Routes
func (rs Routes) ListNames() []string {

	list := make([]string, 0, 5)
	for _, route := range rs {
		list = append(list, route.Name)

	}

	return list
}

// ListURIs provides a list of all recorded route patterns in Routes
func (rs Routes) ListURIs() []string {

	list := make([]string, 0, 5)
	for _, route := range rs {
		list = append(list, route.Pattern)

	}

	return list
}
