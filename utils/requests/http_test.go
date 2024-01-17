package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	headers := map[string]string{}
	_, statusCode, err := Get("https://httpbin.org/get", headers)
	if err != nil {
		log.Println(err.Error())
	}
	if statusCode != http.StatusOK {
		t.Errorf("Get method expected 200 OK, but %d got", statusCode)
	}
}

type Body struct {
	Name string `json:"name"`
}

func TestPost(t *testing.T) {

	body := Body{
		Name: "david",
	}

	bodyByte, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	data, statusCode, err := Post("https://httpbin.org/post", bytes.NewReader(bodyByte))
	if err != nil {
		log.Println(err.Error())
	}
	if statusCode != http.StatusOK {
		t.Errorf("Get method expected 200 OK, but %d got", statusCode)
	}
	fmt.Println(string(data))

	d := map[string]any{}
	err = json.Unmarshal(data, &d)

	if err != nil || d["json"].(map[string]any)["name"].(string) != "david" {
		t.Errorf("Get method expected POST, but %s got", data)
	}
}
