package board

import(
	"github.com/juanfgs/checkers/lib/piece"
	"errors"
	"fmt"
)

// Size of the board
const size = 8


type Board struct {
	Places [][]*piece.Piece
}

func NewBoard() Board {
	board := Board{}
	board.Places = make([][]*piece.Piece, size)
	for i:= range board.Places {
		board.Places[i] = make([]*piece.Piece,size)
	}
	// Place pieces on the board
	board.placeRedPieces()
	board.placeBlackPieces()
	return board
}

// Renders the board as ascii
func (self Board) RenderText() {
	fmt.Println("  a b c d e f g h i j ")
	fmt.Println("  --------------------")
	for rowidx, row := range self.Places {
		fmt.Print(rowidx)
		fmt.Print("|")
		for idx, col := range row {
			if col != nil {
				fmt.Print(col.RenderText())
			} else {
				fmt.Print("0 ")
			}
			if idx == len(row) -1 {
				fmt.Println("|")
			}
		}
	}
	fmt.Println("  --------------------")
}

// Sets tile as selected
func ( self *Board)  SelectTile(x int, y int)  {
	self.deselectTiles()
	if self.Places[x][y] != nil {
		self.Places[x][y].Selected = true
	} 
}

// Deselects all tiles
func (self *Board) deselectTiles(){
	for _,x := range(self.Places) {
		for _,y := range(x) {
			if y != nil{
				y.Selected = false
			}
		}
	}
}

// Places the red pieces on the board
func ( self *Board ) placeRedPieces() {
	for i := 0; i < len(self.Places) ; i++ {
		for j := 0; j < (size / 2) -1; j++ {
			if i % 2 == 0 { // if row is pair
				if j % 2 != 0 { // place on odd columns
					self.Places[i][j] = piece.NewPiece("red")
				}
			} else {
				if j % 2 == 0 { // place on even columns
					self.Places[i][j] = piece.NewPiece("red")
				}
			}
		}

	}
}

//places the black pieces on the board
func ( self *Board ) placeBlackPieces() {
	for i := size -1 ; i > (len(self.Places) / 2)  ; i-- {
		for j := 0 ; j < size    ; j++ {
			if i % 2 == 0 { // if row is pair
				if j % 2 != 0 { // place on odd columns
					self.Places[j][i] = piece.NewPiece("black")
				}
			} else {
				if j % 2 == 0 { // place on even columns
					self.Places[j][i] = piece.NewPiece("black")
				}
			}
		}

	}
}


func (self *Board) MovePiece(x,y,destX,destY int) error {
	if size < destX {
		return errors.New("Illegal move: out of bounds")
	}
	if size < destY {
		return errors.New("Illegal move: out of bounds")
	}
	if self.Places[y][x] == nil {
		return errors.New("Illegal move: no such piece")
	}
	self.Places[destY][destX] = self.Places[y][x]
	self.Places[y][x] = nil
	return nil
}

func (self *Board) MovePieceBottomLeft(x,y int) {
	self.MovePiece(x,y,x-1,y+1)
}

func (self *Board) MovePieceBottomRight(x,y int) {
	self.MovePiece(x,y,x+1,y+1)
}

func (self *Board) MovePieceTopLeft(x,y int) {
	self.MovePiece(y,x,x-1,y-1)
}

func (self *Board) MovePieceTopRight(x,y int) {
	self.MovePiece(y,x,x+1,y-1)
}
