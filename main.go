package main

import (
	"github.com/borogk/hsweeper/ui"
)

func main() {
	u := ui.NewUiWithTitleMenu()
	u.Loop()
}
