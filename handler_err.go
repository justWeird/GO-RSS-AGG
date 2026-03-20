package main

import (
	"net/http"
)

func handlerErr(w http.ResponseWriter, r *http.Request) {
	// This is a simple HTTP handler function that responds to incoming requests.
	// It takes two parameters:
	// - w: an http.ResponseWriter that allows us to write the response back to the client.
	// - r: an *http.Request that contains information about the incoming request.
	// It MUST be defined with this signature to be used as an HTTP handler in Go.

	// call the respondWithError function to send an error response back to the client
	// we pass the http.ResponseWriter, a status code of 500 (Internal Server Error), and an error message
	respondWithError(w, 400, "An error occurred while processing your request") // using 400 (Bad Request) to indicate that the error is due to a client-side issue
}
