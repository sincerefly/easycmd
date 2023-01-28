package requests

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

/*
	发起 Post 请求
*/
func Post(url string, body io.Reader) ([]byte, int, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, 0, err
	}
	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = r.Body.Close()
	}()
	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, r.StatusCode, err
	}
	return resBody, r.StatusCode, nil
}

/*
	发起 Get 请求
*/
func Get(url string, appHeader map[string]string) ([]byte, int, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	for h, v := range appHeader {
		req.Header.Add(h, v)
	}
	var client = &http.Client{
		Timeout: time.Second * 30,
	}
	r, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = r.Body.Close()
	}()
	resBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, r.StatusCode, err
	}
	return resBody, r.StatusCode, nil
}
