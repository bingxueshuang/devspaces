package keyio

import (
	"encoding/hex"
	"errors"
	"github.com/bingxueshuang/devspaces/core"
	"golang.org/x/term"
	"io"
	"os"
)

var ErrNoFile = errors.New("input file not provided")

// ReadKey sets key from hexKey if set,
// otherwise read from filename.
// If filename is empty string.
// Fallback to standard input if stdin is true.
func ReadKey(key core.EllipticKey, filename string, hexKey string, stdin bool) error {
	if hexKey == "" {
		data, err := ReadFile(filename, stdin)
		if err != nil {
			return err
		}
		hexKey = string(data)
	}
	data, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil
	}
	return key.FromBytes(data)
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

// ReadPassword reads a line of input from the terminal
// without local echo and returns the input string.
func ReadPassword() (string, error) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	t := term.NewTerminal(os.Stdin, "")
	return t.ReadPassword("Enter password: ")
}
