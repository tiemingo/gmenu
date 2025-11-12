package gmenu

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Returns the output, error and stderr
func pipeInput(input string, command string, args ...string) (string, error, string) {
	inputData := []byte(input)

	cmd := exec.Command(command, args...)
	cmd.Stdin = bytes.NewReader(inputData)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", err, stderr.String()
	}

	return stdout.String(), nil, ""
}

// Returns RGBA formated as hex
func (c *Color) getHex() string {
	return fmt.Sprintf("%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

func getItemsString(items []string) string {
	itemsS := ""
	for i, item := range items {
		itemsS += string(item)
		if i+1 != len(items) {
			itemsS += "\n"
		}
	}

	return itemsS
}
