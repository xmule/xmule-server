package virgil

import (
	"fmt"
	"net/http"
	"os"
)

// Handler read from disk
// type Handler struct{}

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

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("%s not on disk", filePath)
		return nil
	}

	if blockSize == 0 {
		blockSize = 512000
	}
	block := make([]byte, blockSize)

	bytesRead, err := file.ReadAt(block, offset)
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}

	block = block[:bytesRead]
	fmt.Printf("Accessing file %s: success.", filePath)

	return block
}

func main() {
	http.HandleFunc("/virgil", ServeHTTP)
	err := http.ListenAndServe(":9999", nil)
	fmt.Println(err)
}
