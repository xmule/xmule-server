package charon

import (
	"fmt"
	"net/http"

	"github.com/golang/snappy"

	"github.com/gomodule/redigo/redis"
)

// Handler frontend application that manage authentication and queries redis for in meory data
// otherwise falls back to Homer checking data on storage
type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	//fmt.Fprintf(w, "File read request: %s\n", filePath)

	fmt.Println("Request from ", r.RemoteAddr, "for file: ", filePath)

	result := readBlock(filePath, 0, 0)
	//result := readFile(filePath)
	w.Write(result)
}

// readBlock
func readBlock(filePath string, offset int64, blockSize int64) (result []byte) {
	fmt.Println("Connecting to redis...")
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
		panic(err)
	}
	defer c.Close()
	fmt.Println("Connecting to redis... success!")

	// _, err = c.Do("SET", filePath, snappy.Encode(nil, []byte("ciaoooo")))
	// if err != nil {
	// 	// handle error
	// 	panic(err)
	// }

	reply, err := c.Do("GET", filePath)
	if err != nil {
		// handle error
		panic(err)
	}
	fmt.Printf("Read info from in redis: %s", reply)

	result, err = snappy.Decode(nil, reply.([]byte))
	if err != nil {
		// handle error
		panic(err)
	}

	return result
}

// readFile
func readFile(filePath string) (result []byte) {
	return []byte("NOT IMPLEMENTED")
}

func main() {
	err := http.ListenAndServe(":9999", Handler{})
	fmt.Println(err)
}
