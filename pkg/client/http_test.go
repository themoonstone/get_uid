package client

import (
	"encoding/json"
	"log"

	"testing"
)

func TestSendRequest(t *testing.T) {
	s1 := &Response{}
	client := NewHttpClient("192.168.1.175:8080")
	s, err := client.SendRequest()
	if err != nil {
		log.Fatal(err)
	}
	b, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(b, s1)
	log.Println("id:", s1.String)
}
