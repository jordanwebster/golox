package loxio

import (
	"bytes"
	"errors"
	"io"
)

var Waiting = errors.New("Waiting")

type ChannelMessage struct {
	bytes []byte
}

type ChannelReadWriter struct {
	messages chan ChannelMessage
	buffer   *bytes.Buffer
}

func NewChannelReadWriter() *ChannelReadWriter {
	return &ChannelReadWriter{
		messages: make(chan ChannelMessage),
		buffer:   new(bytes.Buffer),
	}
}

func (rw *ChannelReadWriter) Read(b []byte) (int, error) {
	msg, more := <-rw.messages
	_, err := rw.buffer.Write(msg.bytes)
	if err != nil {
		panic(err)
	}

	n, err := rw.buffer.Read(b)
	if err == io.EOF && more {
		return n, Waiting
	} else {
		return n, err
	}
}

func (rw *ChannelReadWriter) Write(b []byte) (int, error) {
	dest := make([]byte, len(b))
	copy(dest, b)
	rw.messages <- ChannelMessage{bytes: dest}
	return len(b), nil
}

func (rw *ChannelReadWriter) Close() {
	close(rw.messages)
}
