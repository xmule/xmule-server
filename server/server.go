package server

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/golang/protobuf/proto"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func check(e error) {
	if e != nil {
		if e.Error() != "EOF" {
			panic(e)
		}
	}
}

func writeit(f os.File, conn net.Conn, blkSize int, offset int, end chan bool, nThreads int, size int) {
	blk := make([]byte, blkSize)

	//fmt.Println("*****")
	off := offset
	for i := 0; i <= nThreads; i++ {
		if off >= size {
			break
		}
		read, errRead := f.ReadAt(blk, int64(off))
		//fmt.Println("Read bytes:", read)
		data := &MyBlock{
			Block:  blk[:read],
			Offset: int64(off),
		}

		marshData, err := proto.Marshal(data)
		_, err = conn.Write(marshData)
		check(err)
		//fmt.Println("Written bytes:", written)
		if errRead != nil {
			break
		}

		off = off + blkSize*nThreads
	}
	//fmt.Println("-----")
	end <- true
}

// ReadFile send []bytes, offset
func ReadFile(filename string, fout net.Conn) {
	nThreads := 4
	blockSize := 1024

	f, err := os.Open("/home/dciangot/Downloads/chart.png")
	check(err)
	defer f.Close()

	end := make(chan bool)

	stat, err := f.Stat()
	check(err)
	size := stat.Size()

	offset := 0
	for {
		if offset > int(size) {
			break
		}
		if blockSize*nThreads > int(size) {
			nThreads = int(size)/blockSize + 1
			fmt.Println("Reducing threads to:", nThreads)
		}
		for i := 0; i < nThreads; i++ {
			//fmt.Println(offset)
			go writeit(*f, fout, blockSize, offset, end, nThreads, int(size))
			offset = offset + blockSize
		}
		for j := 0; j < nThreads; j++ {
			<-end
		}
	}

}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	// buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	//_, err := conn.Read(buf)
	//if err != nil {
	//	fmt.Println("Error reading:", err.Error())
	//}
	// Send a response back to person contacting us.

	ReadFile("test", conn)
	// Close the connection when you're done with it.
	conn.Close()
}

func listenToWrite(authIp string) {
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	// Listen for an incoming connection.
	for {
		// TODO: timeout
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		} else if conn.RemoteAddr().String() != authIp {
			fmt.Println("Unauthorized")
			continue
		}
		// Handle connections in a new goroutine.
		handleRequest(conn)
	}
}

type helloHandler struct{}

func (h helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, you've hit %s\n", r.URL.Path)
	// Listen for incoming connections.
	// add remoteAddr do the waitinglist
	fmt.Println(r.RemoteAddr)
	go listenToWrite(r.RemoteAddr)

	w.Write([]byte(CONN_HOST + ":" + CONN_PORT))
}

func main() {
	err := http.ListenAndServe(":9999", helloHandler{})
	fmt.Println(err)
}
