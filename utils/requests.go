package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//MakeGetRequest makes GET request for given url with given headers and returns web page contents as bytes with errors if any.
func MakeGetRequest(url string, headers map[string]string) ([]byte, error) {

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	//set the header(s) from header map to request variable
	for headerName, headerValue := range headers {
		request.Header.Set(headerName, headerValue)
	}

	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return bodyBytes, fmt.Errorf("Invalid response code: %d", response.StatusCode)
	}

	if err != nil {
		return nil, err
	}

	return bodyBytes, nil

}
