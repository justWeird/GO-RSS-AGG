package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/justWeird/GO-RSS-AGG/internal/database"
)

// because handlers have specific signatures,
// we need to define the function as a pointer receiver on the dbConfig struct,
// so that we can access the database connection within the handler function.
func (db *dbConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	//decompose the request body into a struct that matches the expected JSON payload for creating a user
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body) // decode the JSON request body into a Go struct. json.NewDecoder creates a new JSON decoder that reads from it.

	params := parameters{}

	err := decoder.Decode(&params) // decode the JSON data into the params struct. The &params syntax is used to pass a pointer to the params variable, allowing the decoder to modify its fields directly.

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err)) // using 400 (Bad Request) to indicate that the error is due to a client-side issue
		return
	}

	// if all goes well so far, we can then call the function.
	//it takes in a context and the parameters required to create a user, which in this case is just the name.
	// since it returns a value, we can assign it to a variable called user, and also check for any errors that may occur during the database operation.
	user, err := db.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),       // generate a new UUID for the user ID using the github.com/google/uuid package
		CreatedAt: time.Now().UTC(), // set the created_at field to the current time
		UpdatedAt: time.Now().UTC(), //set the updated_at field to the current time
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	// call the respondWithJSON function to send a JSON response back to the client
	// we pass the http.ResponseWriter, a status code of 200 (OK), and a simple payload containing a message
	respondWithJSON(w, 200, dbUserToUser(user)) //initialize an empty struct as the payload, which will be converted to an empty JSON object in the response
}
