package main

import (
	"fmt"
	"net/http"
)

func main() {

	fileName := "/tmp/testme6.root"

	url := "http://0.0.0.0:12345/charon?file=" + fileName
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 404 {
		url := "http://0.0.0.0:12345/virgil?file=" + fileName
		resp, err = http.Get(url)
		if err != nil {
			panic(err)
		}
		if resp.StatusCode == 404 {
			url := "http://0.0.0.0:12345/homer?file=" + fileName
			resp, err = http.Get(url)
			if err != nil {
				panic(err)
			}
		}
	}

	body := make([]byte, 512000)
	n, err := resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}

	fmt.Println(string(body[:n]))
}
