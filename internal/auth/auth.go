package auth

import (
	"errors"
	"net/http"
	"strings"
)

/*
GetAPIKey extracts an Api Key from the headers of an HTTP request
Example:
Authorization: Bearer {apikey}
*/
func GetUserByApiKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no Authentication info found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil
}
