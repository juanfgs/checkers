package piece

type Piece struct {
	Level string
	Team string
	Selected bool

}

func NewPiece(team string) *Piece {
	return &Piece{"man",  team, false}
}

func (self *Piece) RenderText() string {
	if self.Team == "red" {
		return "r "
	}
	return "b "
}


