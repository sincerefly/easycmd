package requests

import (
	"errors"
	"github.com/intel-go/fastjson"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

func DoRequest(req *http.Request) (map[string]interface{}, error) {
	if req == nil {
		return nil, errors.New("param error")
	}
	url := req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
	resBody, _, err := Post(url, req.Body)
	if err != nil {
		return nil, err
	}
	var response map[string]interface{}
	err = fastjson.Unmarshal(resBody, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

/*
	向 BytePower 发起 Post 请求
*/
func Post(url string, body io.Reader) ([]byte, int, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, 0, err
	}
	//appHeader := map[string]string{
	//	"Accept":                     config.Conf.App.Accept,
	//	"X-BytePower-Application-Id": config.Conf.App.Id,
	//	"X-BytePower-Auth-Token":     utils.GenerateBPAuthToken(config.Conf.Bp.KeyID, config.Conf.Bp.KeySecret),
	//}
	//for h, v := range appHeader {
	//	req.Header.Add(h, v)
	//}
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
	向 BytePower 发起 Get 请求
*/
func Get(url string, appHeader map[string]string) ([]byte, int, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}
	//appHeader := map[string]string{
	//	"Accept":                     config.Conf.App.Accept,
	//	"X-BytePower-Application-Id": config.Conf.App.Id,
	//	"X-BytePower-Auth-Token":     utils.GenerateBPAuthToken(config.Conf.Bp.KeyID, config.Conf.Bp.KeySecret),
	//}
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
