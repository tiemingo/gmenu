package gmenu

type Color struct {
	R, G, B, A uint8
}

type Gmenu interface {
	GetPrompt() (string, []string)  // Returns command and a slice of args.
	PromptUser() (*string, error)   // Return nil if no item was selected, error if the execution failed and a pointer to an item if it succeeds.
	SetItems(items ...string)       // Sets the items to display.
	AddItems(items ...string)       // Adds new items at the bottom of the list.
	Version() (v string, err error) // Return the version of the menu or an error if the execution failed.
}
