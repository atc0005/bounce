package routes

import (
	"log"
	"net/http"
)

// Route reflects the patterns and handlers for each supported path in our
// API.
type Route struct {
	Name           string
	AllowedMethods []string
	Pattern        string
	Description    string
	HandlerFunc    http.HandlerFunc
}

// Routes is a collection of defined routes, intended for bulk registration
// and auto-indexing on the landing page for this application
type Routes []Route

// Add appends one or many new routes to the existing Routes collection.
func (rs *Routes) Add(r ...Route) {

	for _, newRoute := range r {
		log.Printf("DEBUG: Add %s to routes ...\n", newRoute.Name)
		*rs = append(*rs, newRoute)
	}
}

// RegisterWithServeMux registers each recorded Route with the specified
// HTTP ServeMux.
func (rs *Routes) RegisterWithServeMux(mux *http.ServeMux) {

	// TODO: How would we check for errors registering our route?

	for _, route := range *rs {
		log.Printf("DEBUG: Register %s with ServeMux ...\n", route.Name)
		mux.HandleFunc(route.Pattern, route.HandlerFunc)
	}

}

// ListNames provides a list of all recorded route names in Routes
func (rs Routes) ListNames() []string {

	var list []string
	for _, route := range rs {
		list = append(list, route.Name)

	}

	return list
}

// ListURIs provides a list of all recorded route patterns in Routes
func (rs Routes) ListURIs() []string {

	var list []string
	for _, route := range rs {
		list = append(list, route.Pattern)

	}

	return list
}
