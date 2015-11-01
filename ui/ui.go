package ui

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
	"github.com/juanfgs/checkers/lib/board"
	"log"
)

var RED = []float64{0.8, 0, 0}
var BLACK = []float64{0, 0, 0}

type MainWindow struct {
	*gtk.Window
	MainArea    *gtk.Box
	BoardView   *gtk.DrawingArea
	BoardHeight int
	BoardWidth  int
	boardSize   int
	tileWidth   float64
	tileHeight   float64
	Board       board.Board
	err         error
}

func NewMainWindow() *MainWindow {
	self := new(MainWindow)
	var err error
	self.Window, err = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	if err != nil {
		log.Fatal("Failed to load GTK")
	}

	self.Board = board.NewBoard()


	self.InitializeWidgets()
	self.Window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	self.Window.SetDefaultSize(640, 560)

	return self
}

// Initialize the game widgets
func (self *MainWindow) InitializeWidgets() {

	self.BoardView, self.err = gtk.DrawingAreaNew()
	self.MainArea, self.err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
	self.BoardWidth = 640
	self.BoardHeight = 480

	if self.err != nil {
		log.Fatal("Failed to draw board")
	}

	self.MainArea.PackStart(self.BoardView, true, true, 10)
	self.Window.Add(self.MainArea)
	self.setBoardSize(8)
	self.BoardView.Connect("draw", self.drawBoard)

	self.Window.SetTitle("Checkers")

}


func (self *MainWindow) setBoardSize(size int) {
	self.boardSize = size
	self.tileWidth = float64(self.BoardWidth / self.boardSize)
	self.tileHeight = float64(self.BoardHeight / self.boardSize)
	self.BoardView.Emit("draw")
}

// Draws board with  checkers pattern
func (self *MainWindow) drawBoard(da *gtk.DrawingArea, cr *cairo.Context) bool {


	for i := 0; i < self.boardSize; i++ {
		for j := 0; j < self.boardSize; j++ {
			x := float64(j) * self.tileWidth 
			y := float64(i) * self.tileHeight

			if (i % 2) == (j % 2) {
				cr.Rectangle(x, y, self.tileWidth, self.tileHeight)
				cr.SetSourceRGB(0.5, 0.3, 0)
				cr.Fill()
			} else {
				cr.Rectangle(x, y, self.tileWidth, self.tileHeight)
				cr.SetSourceRGB(0.2, 0, 0)
				cr.Fill()
			}
		}
	}

	for i, row := range self.Board.Places {
		for j, col := range row {
			if col != nil {
				if col.Team == "red" {
					self.DrawPiece(cr, float64(j), float64(i), RED)
				} else {
					self.DrawPiece(cr, float64(j), float64(i), BLACK)
				}
			}
		}
	}

	return false
}

func (self *MainWindow) interactBoard(cr *cairo.Context) {
	cr.Rectangle(0, 0, self.tileWidth, self.tileHeight)
	cr.SetSourceRGB(0.9, 0.9, 0.9)
	cr.Fill()
}


// convenience function to get the position
func (self MainWindow) getPiecePosition(x, y float64) (posx, posy float64) {
	posx = x * self.tileWidth
	posy = y * self.tileHeight
	return posx, posy
}

// draw the piece on the canvas
func (self *MainWindow) DrawPiece(cr *cairo.Context, x, y float64, color []float64) {
	x, y = self.getPiecePosition(x, y)
	x = x + 40
	y = y + 30
	cr.Arc(x, y, 20, 0, 2*3.14)
	cr.SetSourceRGB(color[0], color[1], color[2])
	cr.Fill()
}
