package keyio

import (
	"encoding/hex"
	"errors"
	"github.com/bingxueshuang/devspaces/core"
	"io"
	"os"
)

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
		data, err := readFile(filename, stdin)
		if err != nil {
			return err
		}
		byteslice = data
	}
	return key.FromBytes(byteslice)
}

// read from file.
// fall back to stdin if fallback is set.
func readFile(source string, fallback bool) ([]byte, error) {
	var reader io.Reader = os.Stdin
	if source != "" {
		file, err := os.Open(source)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		reader = file
	} else if !fallback {
		return nil, errors.New("no secret key file provided")
	}
	return io.ReadAll(reader)
}
