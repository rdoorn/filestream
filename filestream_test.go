package filestream

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	os.Remove("file.test")
	stream, err := New("file.test")
	assert.Nil(t, err)
	writer := stream.Writer()
	writer <- "test1"
	writer <- "test2"
	stream.Close()

	stream2, err := New("file.test")
	assert.Nil(t, err)
	reader := stream2.Reader()
	first := <-reader
	assert.Equal(t, "test1", strings.Split(first, ":")[1])
	stream2.Close()
}

func TestStreamWithFilter(t *testing.T) {
	os.Remove("file.test")
	stream, err := New("file.test")
	assert.Nil(t, err)
	writer := stream.Writer()
	writer <- "test1"
	writer <- "test2"
	writer <- "test3"
	writer <- "test4"
	stream.Close()

	filter := func(s string) bool {
		if strings.Index(s, "test1") < 0 {
			return false
		}
		return true
	}

	stream2, err := New("file.test")
	assert.Nil(t, err)
	reader := stream2.ReaderWithFilter(filter)
	first := <-reader // skip test1, return test2
	<-reader          // return test3
	<-reader          // return test4
	eof := <-reader   // return "eof"
	assert.Equal(t, "test2", strings.Split(first, ":")[1])
	assert.Equal(t, "eof", eof)
	stream2.Close()
}
