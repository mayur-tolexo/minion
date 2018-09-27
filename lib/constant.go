package lib

import "net/http"

const PACKAGE_NAME = "github.com/mayur-tolexo/drift/lib"

//Error Messages
const INVALID_REQUEST_MSG = "Invalid Request, Please provide correct input"

// Error Code. PLEASE Map New Error Code To The HTTP Code Map Below.
const (
	VALIDATE_ERROR = 101 //Primary Validation fail
	NO_ERROR       = 0
)

var (
	//ErrorHTTPCode : Error Code to Http Code map
	ErrorHTTPCode = map[int]int{
		NO_ERROR:       http.StatusOK,
		VALIDATE_ERROR: http.StatusBadRequest,
	}
)
