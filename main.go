package main

import (
	"fmt"
	"net/http"

	"github.com/xmule/xmule-server/charon"
	"github.com/xmule/xmule-server/homer"
	"github.com/xmule/xmule-server/virgil"
)

func main() {

	http.HandleFunc("/charon", charon.ServeHTTP)
	http.HandleFunc("/virgil", virgil.ServeHTTP)
	http.HandleFunc("/homer", homer.ServeHTTP)

	fmt.Println(http.ListenAndServe(":12345", nil))

}
