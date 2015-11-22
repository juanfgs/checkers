package piece

const RED = 0
const BLACK = 1

type Piece struct {
	Level string
	Team int
	Selected bool
	Crowned bool

}

func NewPiece(team int) *Piece {
	return &Piece{"man",  team, false, false}
}

func (self *Piece) RenderText() string {
	if self.Team == RED {
		return "r "
	}
	return "b "
}


