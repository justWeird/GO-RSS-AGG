package main

import (
	"fmt"
	"net/http"

	"github.com/justWeird/GO-RSS-AGG/internal/auth"
	"github.com/justWeird/GO-RSS-AGG/internal/database"
)

// we want to use this middle ware to easily call functions (query user, validate API key, etc)
// that will be used across multiple routes that require authentication.

// define a function type for the handler.
// Although it doesn't conform to normal http.HandlerFunc signature,
// it allows us to pass in the database connection as a parameter, which is necessary for our authentication logic.
type authHandler func(http.ResponseWriter, *http.Request, database.User)

// define a middleware function that takes in an authHandler and returns a standard http.HandlerFunc.
func (db *dbConfig) authMiddleware(next authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// extract the API key from the request headers using the GetAPIKey function defined in the internal/auth package.
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, 401, fmt.Sprintf("Auth Error: %v", err)) // using 401 (Unauthorized) to indicate that the error is due to authentication issues
			return
		}

		// query the database for the user associated with the provided API key using a method defined on the database.Queries struct.
		user, err := db.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, 401, fmt.Sprintf("Auth Error: %v", err))
			return
		}

		// if the API key is valid and a user is found, call the next handler in the chain, passing in the user information as a parameter.
		next(w, r, user)
	}
}
