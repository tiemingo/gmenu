# gmenu
Golang bindings for wmenu.


# Wmenu example
``` golang
package main

import (
	"fmt"
	"slices"

	"github.com/gotk3/gotk3/pango"
	"github.com/tiemingo/gmenu"
)

func main() {
	foodSelection := []string{
		"Pizza", "Burger", "Curry",
	}

	gm := gmenu.NewWmenu(gmenu.WmenuOptions{
		BackgroundColor:          &gmenu.Color{R: 255, G: 180, B: 90, A: 255},
		ForegroundColorSelection: &gmenu.Color{R: 0, G: 0, B: 0, A: 255},
		ForegroundColor:          &gmenu.Color{R: 0, G: 0, B: 0, A: 255},
		Prompt: &gmenu.WmenuPrompt{
			Prompt:          "What's your favourite food?",
			BackgroundColor: &gmenu.Color{R: 255, G: 170, B: 0, A: 255},
		},
		Font:         pango.FontDescriptionFromString("Cantarell Italic Bold 12"),
		UseItemLines: true,
	}, foodSelection...)

	resP, err := gm.PromptUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	if resP == nil || !slices.Contains(foodSelection, *resP) {
		fmt.Println("It seems like the users favourite food isn't on the list.")
		return
	}

	fmt.Printf("The users favourite food is %v\n", *resP)
}
```