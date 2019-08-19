package filestream

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type Stream struct {
	file        *os.File
	streamWrite chan string
	streamRead  chan string
	quit        chan struct{}
}

func New(filename string) (*Stream, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	return &Stream{
		file:        file,
		streamWrite: make(chan string),
		streamRead:  make(chan string),
		quit:        make(chan struct{}),
	}, nil
}

func (s *Stream) Writer() chan string {
	go s.writeHandler()
	return s.streamWrite
}

func (s *Stream) writeHandler() {
	log.Printf("writer starting")
	for {
		select {
		case str := <-s.streamWrite:
			s.file.Write([]byte(fmt.Sprintf("%d:%s\n", time.Now().Unix(), str)))
		case <-s.quit:
			log.Printf("writer exiting")
			return
		}
	}
}

func (s *Stream) ReaderWithFilter(filter ...func(string) bool) <-chan string {
	go s.readHandler(filter...)
	return s.streamRead
}

func (s *Stream) Reader() <-chan string {
	go s.readHandler()
	return s.streamRead
}

func (s *Stream) readHandler(filter ...func(string) bool) {
	log.Printf("reader starting")
	scanner := bufio.NewScanner(s.file)

	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		data := scanner.Text()

		write := true
		for _, f := range filter {
			if f(data) {
				write = false
			}
		}
		if write {
			s.streamRead <- data
		}
	}
	s.streamRead <- "eof"
	log.Printf("reader exiting")
}

func (s *Stream) Close() {
	close(s.quit)
}
