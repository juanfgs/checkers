package main

import(
	"github.com/gotk3/gotk3/gtk"
	"github.com/juanfgs/checkers/ui"

)

func main(){
	gtk.Init(nil)
	window := ui.NewMainWindow()
	window.Window.ShowAll()

	gtk.Main()
}
