package gmenu

import (
	"fmt"
	"os/exec"
	"strings"
)

type Rofi struct {
	command string
	items   []string
	opts    RofiOptions
}

type RofiOptions struct {
	Threads         int       // 0: Autodetect the number of supported hardware threads. 1: Disable threading. 2...n: Specify the maximum number of threads to use in the thread pool.
	Display         string    // Display to show on.
	CaseInsensitive bool      // Rofi matches menu items case insensitively.
	CaseSmart       bool      // Start in case-smart mode behave like vim’s smartcase, which determines case-sensitivity by input. Suppresses CaseInsensitive option.
	PredefinedInput string    // Sets input fields content to this.
	ConfigFilePath  string    // Path to configuration file for rofi.
	CachePath       string    // Path to directory that should be used to place temporary files, like history.
	ScrollMethod    int       // Select the scrolling method. 0: Per page, 1: continuous.
	NormalizeMatch  bool      // Normalize the string before matching, so o will match ö, and é matches e. This is not a perfect implementation, but works. For now, it disables highlighting of the matched part.
	DisableLazyGrab bool      // Disables lazy grab, this forces the keyboard being grabbed before gui is shown.
	DisablePlugins  bool      // Disable plugin loading.
	PluginPath      string    // Specify the directory where rofi should look for plugins.
	Markup          bool      // Use Pango markup to format output wherever possible.
	NormalWindow    bool      // Make rofi react like a normal application window. Useful for scripts like Clerk that are basically an application.
	TransientWindow bool      // Make rofi react like a modal dialog that is transient to the currently focused window. Useful when you use a keyboard shortcut to run and show on the window you are working with.
	StealFocus      bool      // Make rofi steal focus on launch and restore close to window that held it when launched.
	Matching        []string  // Specify the matching algorithm used. Multiple matching methods can be specified. The matching up/down keybinding allows cycling through at runtime. Options: normal, regex, glob, fuzzy, prefix
	Tokenize        bool      // Tokenize the input.
	NoSorting       bool      // Enable, disable sort for filtered menu.
	SortingMethod   string    // Sorting options based on input. Methods: levenshtein, fzf. By default levenshtein is used.
	Theme           RofiTheme // Rofi theme.
	Dpi             int       // Override the default DPI setting. If set to 0/1 dpi is set based on auto-detect. For more information see: https://man.archlinux.org/man/rofi.1#dpi
	SelectedRow     int       // Select a certain row. Default: 0.
	Pid             string    // Make rofi create a pid file and check this on startup. The pid file prevents multiple rofi instances from running simultaneously. This is useful when running rofi from a key-binding daemon.
	Replace         bool      // If rofi is already running, based on pid file, try to kill that instance.
	NoClickToExit   bool      // Click the mouse outside the rofi window to exit.
	CustomArgs      []string  // Custom arguments if more are needed.
}

type RofiTheme struct {
	Path  string // Path to theme file, or name of an installed theme. If this option is used the other theme options are ignored. See rofi-theme(5) manpage on how themes are resolved.
	Theme string // Rofi theme in string form. For information on theming check: https://man.archlinux.org/man/rofi-theme.5.en.
}

const (
	MatchingNormal = "normal"
	MatchingRegex  = "regex"
	MatchingGlob   = "glob"
	MatchingFuzzy  = "fuzzy"
	MatchingPrefix = "prefix"

	SortingMethodLevenshtein = "levenshtein"
	SortingMethodFZF         = "fzf"
)

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

	// Set to dmenu mode
	args = append(args, "-dmenu")

	// Threads used
	if r.opts.Threads != 0 {
		args = append(args, "-threads")
		args = append(args, fmt.Sprint(r.opts.Threads))
	}

	// Display
	if r.opts.Display != "" {
		args = append(args, "-display")
		args = append(args, r.opts.Display)
	}

	// Case sensitivity
	if !r.opts.CaseInsensitive {
		args = append(args, "-case-sensitive")
	}

	// Smart case
	if r.opts.CaseSmart {
		args = append(args, "-case-smart")
	}

	// Predefined input
	if r.opts.PredefinedInput != "" {
		args = append(args, "-filter")
		args = append(args, r.opts.PredefinedInput)
	}

	// Config file path
	if r.opts.ConfigFilePath != "" {
		args = append(args, "-config")
		args = append(args, r.opts.ConfigFilePath)
	}

	// Cache path
	if r.opts.CachePath != "" {
		args = append(args, "-cache-dir")
		args = append(args, r.opts.CachePath)
	}

	// Scroll method
	args = append(args, "-scroll-method")
	args = append(args, fmt.Sprint(r.opts.ScrollMethod))

	// Cache path
	if r.opts.NormalizeMatch {
		args = append(args, "-normalize-match")
	}

	// Disable lazy grab
	if r.opts.DisableLazyGrab {
		args = append(args, "-no-lazy-grab")
	}

	// Plugins
	if r.opts.DisablePlugins {
		args = append(args, "-no-plugins")
	}
	if r.opts.PluginPath != "" {
		args = append(args, "-plugin-path")
		args = append(args, r.opts.PluginPath)
	}

	// Markup
	if r.opts.Markup {
		args = append(args, "-markup")
	}

	// Window
	if r.opts.NormalWindow {
		args = append(args, "-normal-window")
	}
	if r.opts.TransientWindow {
		args = append(args, "-transient-window")
	}

	// Steal focus
	if r.opts.StealFocus {
		args = append(args, "-steal-focus")
	} else {
		args = append(args, "-no-steal-focus")
	}

	// Matching
	if len(r.opts.Matching) != 0 {
		args = append(args, "-matching")
		args = append(args, strings.Join(r.opts.Matching, ","))
	}

	// Tokenize
	if r.opts.Tokenize {
		args = append(args, "-tokenize")
	}

	// Sorting
	if r.opts.NoSorting {
		args = append(args, "-no-sort")

	} else {
		args = append(args, "-sort")
	}
	if r.opts.SortingMethod != "" {
		args = append(args, "-sorting-method")
		args = append(args, r.opts.SortingMethod)
	}

	// Theme
	if r.opts.Theme.Path != "" {
		args = append(args, "-theme")
		args = append(args, r.opts.Theme.Path)
	} else if r.opts.Theme.Theme != "" {
		args = append(args, "-theme-str")
		args = append(args, r.opts.Theme.Theme)
	}

	// Dpi
	if r.opts.Dpi != 0 {
		args = append(args, "-dpi")
		args = append(args, fmt.Sprint(r.opts.Dpi))
	}

	// Selected row
	if r.opts.SelectedRow != 0 {
		args = append(args, "-selected-row")
		args = append(args, fmt.Sprint(r.opts.SelectedRow))
	}

	// Pid
	if r.opts.Pid != "" {
		args = append(args, "-pid")
		args = append(args, r.opts.Pid)
	}
	if r.opts.Replace {
		args = append(args, "-replace")
	}

	// No click to exit
	if r.opts.NoClickToExit {
		args = append(args, "-no-click-to-exit")
	} else {
		args = append(args, "-click-to-exit")
	}

	// Add custom arguments
	args = append(args, r.opts.CustomArgs...)

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
