package filestream

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStream(t *testing.T) {
	stream := New("file.test")
	writer := stream.Writer()
	writer <- "test1"
	writer <- "test2"
	stream.Close()

	stream2 := New("file.test")
	reader := stream2.Reader()
	first := <-reader
	assert.Equal(t, "test1", strings.Split(first, ":")[1])
	stream2.Close()
}
