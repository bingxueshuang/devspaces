package keyio

import (
	"fmt"
	"io"
	"os"
)

// WriteFile writes hexadecimal of data to file.
// Fallback to standard output if stdout is set.
func WriteFile(data []byte, filename string, stdout bool) error {
	var writer io.WriteCloser = os.Stdout
	if !stdout {
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		writer = file
	}
	_, err := fmt.Fprintf(writer, "%x", data)
	if err != nil {
		return err
	}
	return nil
}
