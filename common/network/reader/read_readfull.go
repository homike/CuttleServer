package reader

import (
	"io"
	"strconv"
)

type NoBuf struct {
	reader io.Reader
}

func NewNoBuf(r io.Reader) *NoBuf {
	b := &NoBuf{
		reader: r,
	}
	return b
}

func (b *NoBuf) Read() ([]byte, error) {
	header := make([]byte, HEADER_SIZE)
	_, err := io.ReadFull(b.reader, header)
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseInt(string(header), 10, 16)
	if err != nil {
		return nil, err
	}

	data := make([]byte, size)
	_, err = io.ReadFull(b.reader, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
