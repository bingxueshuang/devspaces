package keyio

import (
	"encoding/hex"
	"github.com/bingxueshuang/devspaces/core"
	"io"
	"os"
)

func ReadKey(key core.EllipticKey, filename string, hexKey string) error {
	var byteslice []byte
	if hexKey != "" {
		data, err := hex.DecodeString(hexKey)
		if err != nil {
			return nil
		}
		byteslice = data
	} else {
		data, err := readFile(filename)
		if err != nil {
			return err
		}
		byteslice = data
	}
	return key.FromBytes(byteslice)
}

// read from file, fallback to stdin
func readFile(source string) ([]byte, error) {
	var reader io.Reader = os.Stdin
	if source != "" {
		file, err := os.Open(source)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	}
	return io.ReadAll(reader)
}
