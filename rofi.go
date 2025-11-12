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
	Threads         int    // 0: Autodetect the number of supported hardware threads. 1: Disable threading. 2...n: Specify the maximum number of threads to use in the thread pool.
	Diplay          string // Display to show on.
	CaseInsensitive bool   // Rofi matches menu items case insensitively.
	CaseSmart       bool   // Start in case-smart mode behave like vim’s smartcase, which determines case-sensitivity by input. Suppresses CaseInsensitive option.
	PredefinedInput string // Sets input fields content to this.
	ConfigFilePath  string // Path to configuration file for rofi.
	CachePath       string // Path to directory that should be used to place temporary files, like history.
	ScrollMethod    string // Select the scrolling method. 0: Per page, 1: continuous.
	NormalizeMatch  bool   // Normalize the string before matching, so o will match ö, and é matches e. This is not a perfect implementation, but works. For now, it disables highlighting of the matched part.
	DisableLazyGrap bool   // Disables lazy grab, this forces the keyboard being grabbed before gui is shown.
	DisablePlugins  bool   // Disable plugin loading.
	PluginPath      string // Specify the directory where rofi should look for plugins.
	Markup          bool   // Use Pango markup to format output wherever possible.
	NormalWindow    bool   // Make rofi react like a normal application window. Useful for scripts like Clerk that are basically an application.
	TransientWindow bool   // Make rofi react like a modal dialog that is transient to the currently focused window. Useful when you use a keyboard shortcut to run and show on the window you are working with.
	StealFocus      bool   //
}

// Rofi is only supported in demnu mode
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
	args := []string{}

	return r.command, args
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
