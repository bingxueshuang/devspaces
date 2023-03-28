package keyio

import (
	"fmt"
	"os"
)

// WriteString writes given string to file.
// Fallback to standard output if stdout is set.
func WriteString(data string, filename string, stdout bool) error {
	if filename == "" && stdout {
		fmt.Print(data)
		return nil
	}
	return os.WriteFile(filename, []byte(data), 0644)
}
