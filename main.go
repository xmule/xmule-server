package main

import (
	"fmt"
	"net/http"

	"github.com/dciangot/go-cache/charon"
)

func main() {

	err := http.ListenAndServe(":9999", charon.Handler{})
	fmt.Println(err)

}
