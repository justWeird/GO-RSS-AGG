package main

import (
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// This is a simple HTTP handler function that responds to incoming requests.
	// It takes two parameters:
	// - w: an http.ResponseWriter that allows us to write the response back to the client.
	// - r: an *http.Request that contains information about the incoming request.
	// It MUST be defined with this signature to be used as an HTTP handler in Go.

	// call the respondWithJSON function to send a JSON response back to the client
	// we pass the http.ResponseWriter, a status code of 200 (OK), and a simple payload containing a message
	respondWithJSON(w, 200, struct{}{}) //initialize an empty struct as the payload, which will be converted to an empty JSON object in the response
}
