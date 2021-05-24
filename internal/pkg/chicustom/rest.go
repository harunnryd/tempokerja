package chicustom

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func (handler Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := handler(w, r)
	if err != nil {
		fmt.Println("RESPONSE", data, err)

		WriteError(w, r, err)
		return
	}

	WriteSuccess(w, r, data)
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	response := NewFailed(
		"0001",
		NewResponseDesc("tidak diketahui.", "unknow error."),
		http.StatusInternalServerError,
		NewMeta("1.0", "healthy", "development"),
	)

	compose(w, r, response, response.HttpStatus)
}

func WriteSuccess(w http.ResponseWriter, r *http.Request, data interface{}) {
	response := NewSucces(
		"0000",
		NewResponseDesc("-", "-"),
		data,
		http.StatusOK,
		NewMeta("1.0", "healthy", "development"),
	)

	compose(w, r, response, response.HttpStatus)
}

func compose(w http.ResponseWriter, r *http.Request, response interface{}, httpStatus int) {
	res, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Failed to unmarshal"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, _ = w.Write(res)
}

type Rest interface {
	GetHTTPMethod() string
	GetPattern() string
	GetHandler() Handler
}

type rest struct {
	HTTPMethod string
	Pattern    string
	Handler    Handler
}

func NewRest(httpMethod string, pattern string, handler Handler) Rest {
	return &rest{
		HTTPMethod: httpMethod,
		Pattern:    pattern,
		Handler:    handler,
	}
}

func (rest *rest) GetHandler() Handler {
	return rest.Handler
}

func (rest *rest) GetHTTPMethod() string {
	return rest.HTTPMethod
}

func (rest *rest) GetPattern() string {
	return rest.Pattern
}
