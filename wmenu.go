package gmenu

import (
	"fmt"
	"slices"
	"strings"

	"github.com/gotk3/gotk3/pango"
)

type Wmenu struct {
	command string
	items   []string
	opts    WmenuOptions
}

type WmenuPrompt struct {
	Prompt          string // Defines the prompt to be displayed to the left of the input field.
	BackgroundColor *Color // Defines the prompt background color.
	ForegroundColor *Color // Defines the prompt foreground color.
}

type WmenuOptions struct {
	Bottom                   bool                   // Wmenu appears at the bottom of the screen.
	CaseInsensitive          bool                   // Wmenu matches menu items case insensitively.
	PasswordMode             bool                   // Wmenu will not directly display the keyboard input, but instead replace it with asterisks.
	Font                     *pango.FontDescription // Defines the font used. For more information, see https://docs.gtk.org/Pango/type_func.FontDescription.from_string.html
	Lines                    int                    // Wmenu lists items vertically, with the given number of lines.
	MaxLines                 int                    // Maximum amount of lines diplayed. If 0 or lower this option gets ignored.
	UseItemLines             bool                   // Whether to use the amount of items to set number of lines.
	Output                   string                 // Wmenu is displayed on the output with the given name. For example: eDP-1
	Prompt                   *WmenuPrompt           // Display a promt with styling.
	BackgroundColor          *Color                 // Defines the normal background color.
	ForegroundColor          *Color                 // Defines the normal foreground color.
	BackgroundColorSelection *Color                 // Defines the selection background color.
	ForegroundColorSelection *Color                 // Defines the selection foreground color.
	CustomArgs               []string               // Custom arguments if more are needed
}

func NewWmenu(opts WmenuOptions, items ...string) Wmenu {
	return Wmenu{
		command: "wmenu",
		items:   items,
		opts:    opts,
	}
}

// AddItems implements Gmenu.
func (w *Wmenu) AddItems(items ...string) {
	w.items = append(w.items, items...)
}

// GetPrompt implements Gmenu.
func (w *Wmenu) GetPrompt() (string, []string) {
	args := []string{}

	// Display location
	if w.opts.Bottom {
		args = append(args, "-b")
	}

	// Case sensitivity
	if w.opts.CaseInsensitive {
		args = append(args, "-i")
	}

	// Password mode
	if w.opts.PasswordMode {
		args = append(args, "-P")
	}

	// Font
	if w.opts.Font != nil {
		args = append(args, fmt.Sprint("-f ", w.opts.Font.ToString()))
	}

	// Lines
	if w.opts.UseItemLines {
		w.opts.Lines = len(w.items)
	}
	if w.opts.MaxLines > 0 && w.opts.Lines > w.opts.MaxLines {
		w.opts.Lines = w.opts.MaxLines
	}
	args = append(args, fmt.Sprintf("-l %v", w.opts.Lines))

	// Output
	if w.opts.Output != "" {
		args = append(args, fmt.Sprint("-o ", w.opts.Output))
	}

	// Prompt
	if w.opts.Prompt != nil {
		if w.opts.Prompt.BackgroundColor != nil {
			args = append(args, fmt.Sprintf("-M \"%v\"", w.opts.Prompt.BackgroundColor.getHex()))
		}
		if w.opts.Prompt.ForegroundColor != nil {
			args = append(args, fmt.Sprintf("-m \"%v\"", w.opts.Prompt.ForegroundColor.getHex()))
		}
		if w.opts.Prompt.Prompt != "" {
			args = append(args, fmt.Sprint("-p ", w.opts.Prompt.Prompt))
		}
	}

	// Colors
	if w.opts.BackgroundColor != nil {
		args = append(args, fmt.Sprintf("-N \"%v\"", w.opts.BackgroundColor.getHex()))
	}
	if w.opts.ForegroundColor != nil {
		args = append(args, fmt.Sprintf("-n \"%v\"", w.opts.ForegroundColor.getHex()))
	}
	if w.opts.BackgroundColorSelection != nil {
		args = append(args, fmt.Sprintf("-S \"%v\"", w.opts.BackgroundColorSelection.getHex()))
	}
	if w.opts.ForegroundColorSelection != nil {
		args = append(args, fmt.Sprintf("-s \"%v\"", w.opts.ForegroundColorSelection.getHex()))
	}

	args = append(args, w.opts.CustomArgs...)

	return w.command, args
}

// PromptUser implements Gmenu.
func (w *Wmenu) PromptUser() (*string, error) {
	items := ""
	for i, item := range w.items {
		items += string(item)
		if i+1 != len(w.items) {
			items += "\n"
		}
	}
	_, args := w.GetPrompt()
	outS, err := pipeInput(items, w.command, args...)
	if err != nil {
		return nil, err
	}

	item := strings.TrimSuffix(outS, "\n")
	if slices.Contains(w.items, item) {
		return &item, nil
	}
	return nil, nil

}

// SetItems implements Gmenu.
func (w *Wmenu) SetItems(items ...string) {
	w.items = items
}

// Version implements Gmenu.
func (w *Wmenu) Version() (string, error) {
	return pipeInput("", w.command, "-v")
}
