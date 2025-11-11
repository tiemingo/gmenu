package gmenu

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Color struct {
	R, G, B, A uint8
}

func (c *Color) getHex() string {
	return fmt.Sprintf("%02x%02x%02x%02x", c.R, c.G, c.B, c.A)
}

type Gmenu interface {
	GetPrompt() (string, []string) // Returns command and a slice of args.
	PromptUser() (*string, error)  // Return nil if no item was selected, error if the execution failed and a pointer to an item if it succeeds.
	SetItems(items ...string)      // Sets the items to display.
	AddItems(items ...string)      // Adds new items at the bottom of the list.
	Version() (string, error)      // Return the version of the menu or an error if the execution failed.
}

func pipeInput(input string, command string, args ...string) (string, error) {
	inputData := []byte(input)

	cmd := exec.Command(command, args...)
	cmd.Stdin = bytes.NewReader(inputData)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command execution failed: %v\nStderr: %s", err, stderr.String())
	}

	return stdout.String(), nil
}
