package reader

import (
	"bytes"
	"io"
	"strconv"
)

type ByteBuf struct {
	reader  io.Reader
	readbuf []byte
	databuf *bytes.Buffer
	size    int // last packet length
}

func NewByteBuf(r io.Reader) *ByteBuf {
	b := &ByteBuf{
		reader:  r,
		databuf: bytes.NewBuffer(nil),
		readbuf: make([]byte, 1024),
		size:    -1,
	}
	return b
}

func (b *ByteBuf) header() error {
	header := b.databuf.Next(HEADER_SIZE)
	size, err := strconv.Atoi(string(header))
	if err != nil {
		return err
	}
	b.size = size
	return nil
}

func (b *ByteBuf) Reset() {
	b.databuf = bytes.NewBuffer(nil)
	b.size = -1
}

func (b *ByteBuf) Read() ([][]byte, error) {
	n, err := b.reader.Read(b.readbuf)
	if err != nil {
		return nil, err
	}

	b.databuf.Write(b.readbuf[:n])

	// check length
	if b.databuf.Len() < HEADER_SIZE {
		return nil, err
	}

	// first time
	if b.size < 0 {
		if err = b.header(); err != nil {
			return nil, err
		}
	}

	arrData := [][]byte{}
	for b.size <= b.databuf.Len() {
		arrData = append(arrData, b.databuf.Next(b.size))

		// more packet
		if b.databuf.Len() < HEADER_SIZE {
			b.size = -1
			break
		}

		if err = b.header(); err != nil {
			return nil, err
		}
	}

	return arrData, nil
}
