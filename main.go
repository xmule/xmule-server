package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		if e.Error() != "EOF" {
			panic(e)
		}
	}
}

type block struct {
	data   []byte
	offset int64
}

func writeit(f os.File, fout os.File, blkSize int, offset int, end chan bool, nThreads int, size int) {
	blk := make([]byte, blkSize)

	//fmt.Println("*****")
	off := offset
	for i := 0; i <= nThreads; i++ {
		if off >= size {
			break
		}
		read, errRead := f.ReadAt(blk, int64(off))
		//fmt.Println("Read bytes:", read)

		_, err := fout.WriteAt(blk[:read], int64(off))
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

func main() {
	nThreads := 4
	blockSize := 1024

	f, err := os.Open("/home/dciangot/Downloads/chart.png")
	check(err)
	defer f.Close()

	fout, err := os.Create("/home/dciangot/Downloads/dat_out.png")
	check(err)
	defer fout.Close()

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
			go writeit(*f, *fout, blockSize, offset, end, nThreads, int(size))
			offset = offset + blockSize
		}
		for j := 0; j < nThreads; j++ {
			<-end
		}
	}

}
