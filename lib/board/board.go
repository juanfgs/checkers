package board

import (
	"errors"
	"fmt"
	"github.com/juanfgs/checkers/lib/piece"

	"log"
)

// Size of the board
const size = 8
const RED = 0
const BLACK = 1

type Board struct {
	Places [][]*piece.Piece
	Turn int
}

func NewBoard() Board {
	board := Board{}
	board.Places = make([][]*piece.Piece, size)
	for i := range board.Places {
		board.Places[i] = make([]*piece.Piece, size)
	}
	board.Turn = RED
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
			if idx == len(row)-1 {
				fmt.Println("|")
			}
		}
	}
	fmt.Println("  --------------------")
}

// Sets tile as selected
func (self *Board) SelectTile(x int, y int) bool {
	if self.Places[y][x] != nil {
		self.deselectTiles()
		self.Places[y][x].Selected = true
		return true
	} else {
		return false
	}
}

// Deselects all tiles
func (self *Board) deselectTiles() {
	for _, x := range self.Places {
		for _, y := range x {
			if y != nil {
				y.Selected = false
			}
		}
	}
}

// Places the red pieces on the board
func (self *Board) placeRedPieces() {
	for i := 0; i < len(self.Places); i++ {
		for j := 0; j < (size/2)-1; j++ {
			if i%2 == 0 { // if row is pair
				if j%2 != 0 { // place on odd columns
					self.Places[i][j] = piece.NewPiece(RED)
				}
			} else {
				if j%2 == 0 { // place on even columns
					self.Places[i][j] = piece.NewPiece(RED)
				}
			}
		}

	}
}

//places the black pieces on the board
func (self *Board) placeBlackPieces() {
	for i := size - 1; i > (len(self.Places) / 2); i-- {
		for j := 0; j < size; j++ {
			if i%2 == 0 { // if row is pair
				if j%2 != 0 { // place on odd columns
					self.Places[j][i] = piece.NewPiece(BLACK)
				}
			} else {
				if j%2 == 0 { // place on even columns
					self.Places[j][i] = piece.NewPiece(BLACK)
				}
			}
		}

	}
}

func (self *Board) MovePiece(x, y, destX, destY int) error {

	if size < destX {
		return errors.New("Illegal move: out of bounds")
	}
	if size < destY {
		return errors.New("Illegal move: out of bounds")
	}
	if self.Places[y][x] == nil {
		return errors.New("Illegal move: no such piece")
	}

	if self.Places[y][x].Team != self.Turn  {
		return errors.New("Uh oh: It's not your turn")
	}


	if !self.isLegalMovement(x, y, destX, destY) {
		return errors.New("Illegal move")
	}


	self.Places[destY][destX] = self.Places[y][x]
	self.Places[y][x] = nil



	if self.Turn == RED {
		self.Turn = BLACK
	} else {
		self.Turn = RED
	}

	return nil
}

func (self *Board) MovePieceBottomLeft(x, y int) {
	self.MovePiece(x, y, x-1, y+1)
}

func (self *Board) MovePieceBottomRight(x, y int) {
	self.MovePiece(x, y, x+1, y+1)
}

func (self *Board) MovePieceTopLeft(x, y int) {
	self.MovePiece(y, x, x-1, y-1)
}

func (self *Board) MovePieceTopRight(x, y int) {
	self.MovePiece(y, x, x+1, y-1)
}


// Count the Red pieces
func (self Board) GetScores() (reds,blacks int) {
	reds,blacks = 0, 0
	for _, i := range self.Places {
		for _, j := range i {
			if j != nil && j.Team == RED {
				reds++
			}
			if j != nil && j.Team == BLACK {
				blacks++
			}
		}
	}
	return reds,blacks
}


// Determines wether a movement is legal according to the rules of the Game
func (self Board) isLegalMovement(x, y, destX, destY int) bool {
	if self.Places[y][x] != nil {
		if destX == x || destY == y {
			return false
		}

		if ((destX - x   > 1 || destY - y > 1) || (x - destX   > 1 ||  y - destY > 1))  && !self.Places[y][x].Crowned {
			target := self.pieceBetween(x,y,destX,destY)
			log.Println(target)
			if self.Places[target[1]][target[0]] != nil && target[1] != x && target[0] != y {
				log.Println("eat it!")
				log.Println(self.Places[target[1]][target[0]] )
				self.Places[target[1]][target[0]] = nil
				return true
			}

			
			return false



		}


		return true
	} else {
		return false
	}
}

// Determines if there is a  piece between origin and destination
func (self Board) pieceBetween(x,y,destX,destY int) []int {
	var targetPiece []int = make([]int,2 )

	if destX > x {
		targetPiece[0] =  destX - 1
	} else {
		targetPiece[0] =  destX + 1
	}

	if destY > y {
		targetPiece[1] = destY - 1
	} else {
		targetPiece[1] = destY + 1
	}

	return targetPiece
}
