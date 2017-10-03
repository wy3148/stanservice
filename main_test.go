package main

import (
	"fmt"
	"testing"
)

func TestFilterShow(t *testing.T) {

	b, err := filterShow(nil)

	//empty array in request, return empty array in response
	if err != nil {
		t.Fatal(err.Error())
	}

	if string(b) != `{"response":null}` {
		t.Fatal("wrong response")
	}

	var r ShowRequest
	r.Payload = nil

	//empty payload return empty array
	b, err = filterShow(&r)

	if err != nil {
		t.Fatal(err.Error())
	}

	if string(b) != `{"response":null}` {
		t.Fatal("wrong response")
	}

	items := make([]Show, 2)
	items[0].Drm = false
	items[0].EpisodeCount = 1000

	items[1].Drm = true
	items[1].EpisodeCount = 666

	r.Payload = items

	b, err = filterShow(&r)

	if err != nil {
		t.Fatal(err.Error())
	}

	items[1].Drm = false

	fmt.Printf("one show return:%s\n", string(b))

	items[1].Drm = false

	b, err = filterShow(&r)

	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("return zero show return:%s\n", string(b))
}
