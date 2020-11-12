package reader

import (
	"log"
	"os"
	"syscall"
	"testing"
)

func BenchmarkRead_BufIO(b *testing.B) {
	b.StopTimer()

	f, err := os.Open("data")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := NewBufio(f)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := buf.Read()
		if err != nil {
			panic(err)
		}
		_, err = buf.Read()
		if err != nil {
			panic(err)
		}

		buf.Reader.Reset(f)
		syscall.Seek(int(f.Fd()), 0, 0)
	}
}

func BenchmarkRead_ByteBuf(b *testing.B) {
	b.StopTimer()

	f, err := os.Open("data")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := NewByteBuf(f)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, err = buf.Read()
		if err != nil {
			panic(err)
		}

		buf.Reset()
		syscall.Seek(int(f.Fd()), 0, 0)
	}
}

func BenchmarkRead_ReadFull(b *testing.B) {
	b.StopTimer()

	f, err := os.Open("data")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	buf := NewNoBuf(f)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := buf.Read()
		if err != nil {
			panic(err)
		}
		_, err = buf.Read()
		if err != nil {
			panic(err)
		}
		syscall.Seek(int(f.Fd()), 0, 0)
	}
}

func _TestRead(t *testing.T) {
	{
		f, err := os.Open("data")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		buf := NewBufio(f)
		data, err := buf.Read()
		if err != nil {
			panic(err)
		}
		log.Println("length:", len(data), "message:", string(data))

		data, err = buf.Read()
		if err != nil {
			panic(err)
		}
		log.Println("length:", len(data), "message:", string(data))
	}

	{
		f, err := os.Open("data")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		buf := NewByteBuf(f)
		datas, err := buf.Read()
		if err != nil {
			panic(err)
		}
		for _, v := range datas {
			log.Println("length:", len(v), "message:", string(v))
		}
	}

	// copy once, call sys.read twice
	{
		f, err := os.Open("data")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		buf := NewNoBuf(f)
		data, err := buf.Read()
		if err != nil {
			panic(err)
		}
		log.Println("length:", len(data), "message:", string(data))

		data, err = buf.Read()
		if err != nil {
			panic(err)
		}
		log.Println("length:", len(data), "message:", string(data))
	}
}
