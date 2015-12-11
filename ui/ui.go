package ui

import (
	"fmt"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/juanfgs/checkers/lib/board"
	"log"
	"errors"
)

var RED = []float64{0.8, 0, 0}
var BLACK = []float64{0, 0, 0}

type MainWindow struct {
	*gtk.Window
	// Areas
	MainArea  *gtk.Box
	ScoreArea *gtk.Box

	// Labels
	Scores      *gtk.Label
	BlacksScore *gtk.Label
	RedsScore   *gtk.Label

	// Board drawing details
	BoardEventBox *gtk.EventBox
	BoardView     *gtk.DrawingArea
	BoardHeight   int
	BoardWidth    int
	boardSize     int
	tileWidth     float64
	tileHeight    float64

	// Internal board structure
	Board board.Board
	err   error
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
	self.MainArea, self.err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 1)
	self.ScoreArea, self.err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 2)
	self.BoardEventBox, self.err = gtk.EventBoxNew()
	self.Scores, self.err = gtk.LabelNew("Scores:")
	self.BlacksScore, self.err = gtk.LabelNew("Black pieces:")
	self.RedsScore, self.err = gtk.LabelNew("Black pieces:")
	self.BoardWidth = 640
	self.BoardHeight = 480
	self.BoardEventBox.Add(self.BoardView)
	if self.err != nil {
		log.Fatal("Failed to draw board")
	}

	self.BoardEventBox.AddEvents(int(gdk.BUTTON_PRESS_MASK))
	self.BoardEventBox.Connect("button_press_event", self.interactBoard)
	self.MainArea.PackStart(self.BoardEventBox, true, true, 10)
	self.MainArea.Add(self.ScoreArea)

	self.Scores.SetMarkup("<strong>%s</strong>")
	self.Scores.SetJustify(gtk.JUSTIFY_CENTER)

	self.ScoreArea.Add(self.Scores)
	self.BlacksScore.SetJustify(gtk.JUSTIFY_LEFT)
	self.RedsScore.SetJustify(gtk.JUSTIFY_LEFT)
	self.ScoreArea.Add(self.BlacksScore)
	self.ScoreArea.Add(self.RedsScore)

	self.Window.Add(self.MainArea)
	self.setBoardSize(8)
	self.BoardView.Connect("draw", self.drawBoard)
	self.Board.RenderText()

	self.Window.SetTitle("Checkers")
}

func (self *MainWindow) drawScores() {
	reds, blacks := self.Board.GetScores()

	self.BlacksScore.SetText(fmt.Sprintf("Black Pieces: %d", blacks))
	self.RedsScore.SetText(fmt.Sprintf("Red Pieces: %d", reds))
}

func (self *MainWindow) setBoardSize(size int) {
	self.boardSize = size
	self.tileWidth = float64(self.BoardWidth / self.boardSize)
	self.tileHeight = float64(self.BoardHeight / self.boardSize)
	self.BoardView.Emit("draw")
}

// Draws and re-draws the board
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
	// Draw pieces
	for i, row := range self.Board.Places {
		for j, col := range row {
			if col != nil {
				if col.Selected {
					self.DrawSelector(cr, float64(i), float64(j))
				}
				if col.Team == board.RED {
					self.DrawPiece(cr, float64(i), float64(j), RED)
				} else {
					self.DrawPiece(cr, float64(i), float64(j), BLACK)
				}
			}
		}
	}

	self.drawScores()
	return false
}

// Handles user interaction with the board
func (self *MainWindow) interactBoard(eb *gtk.EventBox, event *gdk.Event) {
	evbutton := &gdk.EventButton{event}

	y, x, error := self.calculatePosition(evbutton.X(), evbutton.Y())

	if  error == nil && !self.Board.SelectTile(x, y) {
		lastx, lasty := self.getSelectedPiece()
		if lastx != -1 && self.Board.Places[y][x] == nil {
			err := self.Board.MovePiece(lastx, lasty, x, y)

			if err != nil {
				errorMessage, _ := gtk.LabelNew(err.Error())
				alert, _ := gtk.DialogNew()
				alert.SetTitle("Error")
				alert.AddButton("Dismiss", gtk.RESPONSE_CLOSE)
				box, _ := alert.GetContentArea()
				box.Add(errorMessage)
				box.ShowAll()
				if alert.Run() == int(gtk.RESPONSE_CLOSE) {
					alert.Destroy()
				}

			} else {
				self.Board.Places[y][x].Selected = false
			}

		}
	}

	self.BoardView.QueueDraw()
}

func (self MainWindow) getSelectedPiece() (x, y int) {
	for y, row := range self.Board.Places {
		for x, col := range row {
			if col != nil && col.Selected {
				return x, y
			}
		}
	}
	return -1, -1
}

// Calculates the position in the board based on mouse coordinates
func (self MainWindow) calculatePosition(x float64, y float64) (Y, X int, err error) {
	ctx := 0
	cty := 0

	for i := 0; i < int(x); i = i + int(self.tileWidth) {
		Y = ctx
		ctx++
	}

	for j := 0; j < int(y); j = j + int(self.tileHeight) {

		X = cty
		cty++
	}

	if self.boardSize < Y || self.boardSize < X {
		err = errors.New("out of bounds")
	}

	return Y, X, err
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

// draw the piece on the canvas
func (self *MainWindow) DrawSelector(cr *cairo.Context, x, y float64) {

	cr.SetSourceRGBA(0.2, 0.8, 0.2, 0.8)
	cr.SetLineWidth(5)
	cr.Rectangle(x*self.tileWidth, y*self.tileHeight, self.tileWidth, self.tileHeight)
	cr.Stroke()

}
