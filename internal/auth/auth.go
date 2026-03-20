package auth

import (
	"errors"
	"net/http"
	"strings"
)

// This package will contain all the logic related to authentication and authorization,
// such as generating API keys, validating API keys, and handling user authentication.

// this function will be responsible for extracting the API key from the request headers
// and validating it against the database. If the API key is valid, it will return the
// associated user information; otherwise, it will return an error.
// header format: Authorization: ApiKey {insert API key here}
func GetAPIKey(headers http.Header) (string, error) {

	val := headers.Get("Authorization") // extract the value of the "Authorization" header from the request headers using the Get method of the http.Header type.

	if val == "" {
		return "", errors.New("no authentication info found")
	}

	// since the expected format of the header is "ApiKey {insert API key here}", we can split the header value by space
	// and check if the first part is "ApiKey" and the second part is the actual API key.
	vals := strings.Split(val, " ")

	if len(vals) != 2 || vals[0] != "ApiKey" {

		return "", errors.New("malformed header format.")
	}

	return vals[1], nil
}
