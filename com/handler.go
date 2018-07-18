package com

import (
	"fmt"
	"log"
	"net/http"
)

// Handler http://blog.golang.org/error-handling-and-go
type Handler func(http.ResponseWriter, *http.Request) *Error

func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

// Error by Handler call
type Error struct {
	Error   error
	Message string
	Code    int
}

// Errorf ...
func Errorf(err error, format string, v ...interface{}) *Error {
	return &Error{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
