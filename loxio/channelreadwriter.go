package loxio

import "bytes"

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
	for rw.buffer.Len() < len(b) {
		msg, more := <-rw.messages
		_, err := rw.buffer.Write(msg.bytes)
		if err != nil {
			panic(err)
		}

		if !more {
			break
		}
	}

	return rw.buffer.Read(b)
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
