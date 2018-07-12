package charon

import (
	"fmt"
	"net/http"

	"github.com/gomodule/redigo/redis"
)

// Handler frontend application that manage authentication and queries redis for in meory data
// otherwise falls back to Homer checking data on storage
//type Handler struct{}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	//fmt.Fprintf(w, "File read request: %s\n", filePath)

	fmt.Println("Request from ", r.RemoteAddr, "for file: ", filePath)

	result := readBlock(filePath, 0, 0)
	//result := readFile(filePath)
	if result == nil {
		w.WriteHeader(http.StatusNotFound)
	}
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

	if reply == nil {
		fmt.Printf("%s not in memory", filePath)
		return nil
	}

	fmt.Printf("Read info from in redis: %s", reply)

	//result, err = snappy.Decode(nil, reply.([]byte))
	result = reply.([]byte)
	if err != nil {
		// handle error
		panic(err)
	}

	// TODO: send to Dante for writing with consistent prio
	//       but maybe send only the memory key
	//data := &server.MyBlock{
	//	Block:  reply.([]byte),
	//	Offset: int64(offset),
	//}
	//marshData, err := proto.Marshal(data)
	//_, err = conn.Write(marshData)
	//check(err)
	//if errRead != nil {
	//	break
	//}

	return result
}

// readFile
func readFile(filePath string) (result []byte) {
	return []byte("NOT IMPLEMENTED")
	// nThreads := 4
	// blockSize := 1024

	// f, err := os.Open("/home/dciangot/Downloads/chart.png")
	// check(err)
	// defer f.Close()

	// end := make(chan bool)

	// stat, err := f.Stat()
	// check(err)
	// size := stat.Size()

	// offset := 0
	// for {
	// 	if offset > int(size) {
	// 		break
	// 	}
	// 	if blockSize*nThreads > int(size) {
	// 		nThreads = int(size)/blockSize + 1
	// 		fmt.Println("Reducing threads to:", nThreads)
	// 	}
	// 	for i := 0; i < nThreads; i++ {
	// 		//fmt.Println(offset)
	// 		go writeit(*f, fout, blockSize, offset, end, nThreads, int(size))
	// 		offset = offset + blockSize
	// 	}
	// 	for j := 0; j < nThreads; j++ {
	// 		<-end
	// 	}
	// }
}

func main() {
	http.HandleFunc("/charon", ServeHTTP)
	err := http.ListenAndServe(":9999", nil)
	fmt.Println(err)
}
