package gmenu

import (
	"os/exec"
	"strings"
)

type Rofi struct {
	command string
	items   []string
	opts    RofiOptions
}

type RofiOptions struct {
}

func NewRofi(opts RofiOptions, items ...string) Gmenu {
	return &Rofi{
		command: "rofi",
		items:   items,
		opts:    opts,
	}
}

// AddItems implements Gmenu.
func (r *Rofi) AddItems(items ...string) {
	r.items = append(r.items, items...)
}

// GetPrompt implements Gmenu.
func (r *Rofi) GetPrompt() (string, []string) {
	panic("unimplemented")
}

// PromptUser implements Gmenu.
func (r *Rofi) PromptUser() (*string, error) {
	items := getItemsString(r.items)
	_, args := r.GetPrompt()

	outS, err, stderr := pipeInput(items, r.command, args...)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 1 && stderr == "" {
			return nil, nil
		}
		return nil, err
	}

	item := strings.TrimSuffix(outS, "\n")
	return &item, nil
}

// SetItems implements Gmenu.
func (r *Rofi) SetItems(items ...string) {
	r.items = items
}

// Version implements Gmenu.
func (r *Rofi) Version() (v string, err error) {
	v, err, _ = pipeInput("", r.command, "-version")
	return
}
