package web

import (
	"contactApp/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// UnmarshalJSON parses data from request and return otherwise error return.
func UnmarshalJSON(request *http.Request, out interface{}) error {
	if request.Body == nil {
		fmt.Println("==============================err request.Body == nil==========================")
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}
	// fmt.Println("==============================err (request.Body)==========================")
	// fmt.Println(request.Body)
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println("==============================err ioutil.ReadAll==========================")
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}

	if len(body) == 0 {
		fmt.Println("==============================err len(body) == 0==========================")
		return errors.NewHTTPError(errors.ErrorCodeEmptyRequestBody, http.StatusBadRequest)
	}

	err = json.Unmarshal(body, out)
	if err != nil {
		fmt.Println("==============================err json.Unmarshal==========================")
		fmt.Println(body)
		fmt.Println(out)
		return errors.NewHTTPError(err.Error(), http.StatusBadRequest)
	}
	return nil
}
