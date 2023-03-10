package keyio

import (
	"encoding/hex"
	"errors"
	"github.com/bingxueshuang/devspaces/core"
	"io"
	"os"
)

var ErrNoFile = errors.New("input file not provided")

// ReadKey sets key from hexKey if set,
// otherwise read from filename.
// If filename is empty string.
// Fallback to standard input if stdin is true.
func ReadKey(key core.EllipticKey, filename string, hexKey string, stdin bool) error {
	var byteslice []byte
	if hexKey != "" {
		data, err := hex.DecodeString(hexKey)
		if err != nil {
			return nil
		}
		byteslice = data
	} else {
		data, err := ReadFile(filename, stdin)
		if err != nil {
			return err
		}
		byteslice = data
	}
	return key.FromBytes(byteslice)
}

// ReadFile reads bytes from file.
// fall back to stdin if fallback is set.
func ReadFile(source string, fallback bool) ([]byte, error) {
	var reader io.Reader = os.Stdin
	if source != "" {
		file, err := os.Open(source)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	} else if !fallback {
		return nil, ErrNoFile
	}
	return io.ReadAll(reader)
}
