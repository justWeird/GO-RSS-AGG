package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	/*
		This function is responsible for sending an error response in JSON format to the client.
		We are concerend with errors on servier side i.e 500 and above
	*/
	if statusCode > 499 {
		// the error is a server error, so we log the error message for debugging purposes
		log.Printf("Server error 5XX: %s\n", message)
	}
	// we create a payload that contains the error message, which will be sent back to the client as JSON
	// this is defined as an anonymous struct assigned to a variable with a single field "error" that holds the error message
	// the field is tagged with `json:"error"` to specify that when this struct is converted to JSON, the field will be named "error"
	payload := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}
	// we call the respondWithJSON function to send the error response back to the client, passing the payload we just created
	respondWithJSON(w, statusCode, payload)

}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	/*
			This function is responsible for sending a JSON response to the client. It takes three parameters:
		- w: an http.ResponseWriter that allows us to write the response back to the client.
		- statusCode: an integer representing the HTTP status code to be sent in the response.
		- payload: an interface{} that can hold any type of data that we want to send as JSON.
	*/

	// convert the payload to JSON format
	// json.Marshal is a function from the encoding/json package that converts a Go data structure into JSON format.
	// It returns the JSON as a byte slice and an error if the conversion fails.
	response, err := json.Marshal(payload)

	if err != nil {
		// if there is an error during JSON conversion, we set the response status to 500 (Internal Server Error)
		// and write an error message back to the client
		log.Println("Error converting payload to JSON:", payload)
		w.WriteHeader(500)
		return
	}

	// if the conversion is successful, we set the Content-Type header to application/json to indicate that the response is in JSON format
	w.Header().Add("Content-Type", "application/json")
	// write the JSON response back to the client
	w.WriteHeader(statusCode)
	w.Write(response)

}
