package requests

import (
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	headers := map[string]string{}
	data, statusCode, err := Get("https://ifconfig.me/method", headers)
	if err != nil {
		log.Println(err.Error())
	}
	if statusCode != 200 {
		t.Errorf("Get method expected 200 OK, but %d got", statusCode)
	}
	if string(data) != "GET" {
		t.Errorf("Get method expected GET, but %s got", data)
	}
}
