package main

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

const NUMBER_OF_WORKERS_READER int = 5

func main() {
	stream := Simulation()

	reader := New(18)

	reader.Open("./data")

	var wg sync.WaitGroup

	wg.Add(NUMBER_OF_WORKERS_READER)
	for i := 0; i < NUMBER_OF_WORKERS_READER; i++ {
		go Read(i, stream, &wg, reader)
	}

	start := time.Now()
	wg.Wait()

	reader.Close()
	fmt.Println(time.Since(start))
}

func Simulation() <-chan struct{} {
	streaming := make(chan struct{})

	go func() {
		for i := 0; i < 1_000_000; i++ {
			streaming <- struct{}{}
		}
		close(streaming)
	}()

	return streaming
}

func Read(id int, stream <-chan struct{}, wg *sync.WaitGroup, reader *Reader) {
	for range stream {
		reader.Read()
	}

	wg.Done()
}

type Reader struct {
	block []byte
	size  int
	file  *os.File
}

func New(size int) *Reader {
	return &Reader{
		block: make([]byte, size),
		size:  size,
	}
}

func (r *Reader) Open(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	r.file = file
}

func (r *Reader) Close() {
	r.file.Close()
}

func (r *Reader) Read() {
	if r.file == nil {
		fmt.Println("file nije otvoren!")
	}

	io.ReadFull(r.file, r.block)
}
