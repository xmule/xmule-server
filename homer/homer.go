package homer

import (
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	//fmt.Fprintf(w, "File read request: %s\n", filePath)

	fmt.Println("Request from ", r.RemoteAddr, "for file: ", filePath)

	fmt.Println("Connecting to redis...")
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
		panic(err)
	}
	defer c.Close()
	fmt.Println("Connecting to redis... success!")

	result := readBlock(filePath, 0, 0)

	if result == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write(result)
		return
	}
	w.Write(result)

	//_, err = c.Do("SET", filePath, snappy.Encode(nil, result))
	_, err = c.Do("SET", filePath, result)
	if err != nil {
		panic(err)
	}
}

func readBlock(filePath string, offset int64, blockSize int64) (result []byte) {

	url := "http://0.0.0.0:12345/virgil?file=/home/dciangot" + filePath
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body := make([]byte, 512000)
	bytesRead, err := resp.Body.Read(body)
	defer resp.Body.Close()

	if err != nil && err.Error() != "EOF" {
		panic(err)
	}
	fmt.Println(string(body[:bytesRead]))
	return body[:bytesRead]
}

func main() {
	http.HandleFunc("/homer", ServeHTTP)
	err := http.ListenAndServe(":9999", nil)
	fmt.Println(err)
}
