package main

import (
	"bufio"
	"sync"
	"testing"
)

func BenchmarkRead(b *testing.B) {
	b.StopTimer()
	reader := New(4096)
	reader.Open("./data")

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		reader.Read()
	}

	reader.Close()
}

func BenchmarkParallelRead(b *testing.B) {
	b.StopTimer()
	streaming := make(chan struct{})

	reader := New(4096)
	reader.Open("./data")

	var wg sync.WaitGroup
	wg.Add(NUMBER_OF_WORKERS_READER)
	for i := 0; i < NUMBER_OF_WORKERS_READER; i++ {
		go Read(i, streaming, &wg, reader)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		streaming <- struct{}{}
	}
	close(streaming)

	wg.Wait()
	reader.Close()
}

func BenchmarkFastRead(b *testing.B) {
	b.StopTimer()
	reader := New(4096)
	reader.Open("./data")

	bufReader := bufio.NewReader(reader.file)

	buffer := make([]byte, 100)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		bufReader.Read(buffer)
	}
}
