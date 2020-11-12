package reader

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
)

type ByteBuf struct {
	reader  io.Reader
	buf     *bytes.Buffer
	readbuf []byte
	size    int // last packet length
}

func NewByteBuf(r io.Reader) *ByteBuf {
	b := &ByteBuf{
		reader:  r,
		buf:     bytes.NewBuffer(nil),
		readbuf: make([]byte, 1024),
		size:    -1,
	}
	return b
}

func (b *ByteBuf) header() error {
	header := b.buf.Next(HEADER_SIZE)
	size, err := strconv.Atoi(string(header))
	if err != nil {
		return err
	}
	b.size = size
	return nil
}

func (b *ByteBuf) Read() ([][]byte, error) {
	n, err := b.reader.Read(b.readbuf)
	if err != nil {
		log.Println(fmt.Sprintf("Read message error: %v", err))
		return nil, err
	}

	b.buf.Write(b.readbuf[:n])

	// check length
	if b.buf.Len() < HEADER_SIZE {
		return nil, err
	}

	// first time
	if b.size < 0 {
		if err = b.header(); err != nil {
			return nil, err
		}
	}

	arrData := [][]byte{}
	for b.size <= b.buf.Len() {
		arrData = append(arrData, b.buf.Next(b.size))

		// more packet
		if b.buf.Len() < HEADER_SIZE {
			b.size = -1
			break
		}

		if err = b.header(); err != nil {
			return nil, err
		}
	}

	return arrData, nil
}
