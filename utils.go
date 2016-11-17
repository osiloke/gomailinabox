package mailinabox

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

func responseBytes(resp *http.Response) ([]byte, error) {
	return ioutil.ReadAll(resp.Body)
}

func responseHandler(resp *http.Response, err error) (*[]byte, *http.Response, error) {
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, resp, err
	}
	b, _err := responseBytes(resp)
	if err != nil {
		log.Println(_err.Error())
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 500 {
			return nil, resp, ErrServerError
		}
		return nil, resp, errors.New(string(b))
	}
	return &b, resp, nil
}
