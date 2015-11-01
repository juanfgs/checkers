package piece

type Piece struct {
	Level string
	Team string

}

func NewPiece(team string) *Piece {
	return &Piece{"man",  team}
}

func (self *Piece) RenderText() string {
	if self.Team == "red" {
		return "r "
	}
	return "b "
}


