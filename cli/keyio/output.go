package keyio

import (
	"fmt"
	"github.com/bingxueshuang/devspaces/core"
	"io"
	"os"
)

// WriteKey writes to file.
// Fallback to standard output if stdout is set.
func WriteKey(key core.EllipticKey, filename string, stdout bool) error {
	var writer io.WriteCloser = os.Stdout
	if filename != "" {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = file
	}
	if stdout {
		_, err := fmt.Fprintf(writer, "%x", key.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}
