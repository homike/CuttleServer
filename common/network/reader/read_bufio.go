package reader

import (
	"bufio"
	"io"
	"strconv"
)

// proto
//  ----------------
// | size | content |

const (
	HEADER_SIZE = 2 // 2 byte head size
)

type Bufio struct {
	Reader *bufio.Reader
}

func NewBufio(r io.Reader) *Bufio {
	b := &Bufio{}
	b.Reader = bufio.NewReaderSize(r, 1024)

	return b
}

func (this *Bufio) Read() ([]byte, error) {
	var (
		err    error
		buffer []byte
	)

	buffer, err = this.Reader.Peek(HEADER_SIZE)
	if err != nil {
		return nil, err
	}

	size, err := strconv.Atoi(string(buffer))
	if err != nil {
		return nil, err
	}

	totalSize := HEADER_SIZE + size
	buffer, err = this.Reader.Peek(totalSize)
	if err != nil {
		return nil, err
	}

	data := make([]byte, totalSize-HEADER_SIZE)
	copy(data, buffer[HEADER_SIZE:])

	_, err = this.Reader.Discard(totalSize)
	if err != nil {
		return nil, err
	}
	return data, nil
}
