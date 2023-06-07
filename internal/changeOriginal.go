package internal

import (
	"bytes"
	"compress/gzip"
	"io"
	"math/rand"
	"os"
	"time"
)

func GetPEFileData(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func EncodeFile(data []byte) ([]byte, byte) {
	rand.Seed(time.Now().UnixMilli())
	s := (byte)(rand.Intn(255))
	res := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		res[i] = data[i] ^ s
	}
	return res, s
}

func CompressFile(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	w.Flush()
	w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func CopyFile(f1, f2 string) error {
	src, err := os.Open(f1)
	if err != nil {
		return err
	}
	dest, err := os.Create(f2)
	if err != nil {
		return err
	}
	_, err = io.Copy(dest, src)
	if err != nil {
		return err
	}
	dest.Close()
	src.Close()
	return nil
}
