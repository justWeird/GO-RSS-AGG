package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	// create a new router
	router := chi.NewRouter()

	// set up cors routing which allows cross-origin requests from any origin and supports common HTTP methods
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},                   // allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // specify allowed HTTP methods
		AllowedHeaders:   []string{"*"},                                       // specify allowed headers
		ExposedHeaders:   []string{"Link"},                                    // specify headers that can be exposed to the browser
		AllowCredentials: false,                                               // do not allow credentials (cookies, authorization headers, etc.)
		MaxAge:           300,                                                 // maximum age for preflight requests in seconds
	}))

	//  routes are defined to handle incoming HTTP requests. Similar to routes in express which specify how the server should respond to different HTTP methods and paths
	// the handler can be a function that takes an http.ResponseWriter and an *http.Request as parameters, which is the standard signature for HTTP handlers in Go
	v1Router := chi.NewRouter()
	// using HandlerFunc allows the route to accept any type of HTTP request (GET, POST, etc.) and respond accordingly based on the logic defined in the handler function.
	// v1Router.HandleFunc("/health", handler)
	v1Router.Get("/health", handler) // using Get method to specify that this route should only respond to GET requests.
	router.Mount("/v1", v1Router)    // mount the v1Router on the main router at the /v1 path, so that all routes defined in v1Router will be accessible under the /v1 path (e.g., /v1/ready)

	// define an additional route to demonstrate error handling
	v1Router.Get("/err", handlerErr)

	// set up the server
	serverObj := &http.Server{
		Handler: router,     //server requires a handler to route requests
		Addr:    ":" + port, // specify the port to listen on
	}

	fmt.Printf("Server will start on port: %s\n", port)

	err = serverObj.ListenAndServe() // start the server and listen for incoming requests

	// as the server is running, it will block the main goroutine until it receives a shutdown signal or encounters an error
	// log the error if the server fails to start or encounters an issue while running
	if err != nil {
		log.Fatal(err)
	}
}
