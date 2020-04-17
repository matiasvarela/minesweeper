package board

const (
	EMPTY int = 0
	BOMB  int = 1

	STATUS_NEW      string = "new"
	STATUS_ON_GOING string = "on_going"
	STATUS_LOST     string = "lost"
	STATUS_WON      string = "won"
)

type Square struct {
	Type     int  `json:"type"`
	Revealed bool `json:"revealed"`
	Marked   bool `json:"marked"`
}

type SquarePosition struct {
	Row    int `json:"row"`
	Column int `json:"column"`
}

type Board struct {
	Squares              [][]Square        `json:"squares"`
	BombsNumber          int               `json:"bombs_number,omitempty"`
	BombsPositions       *[]SquarePosition `json:"bombs_positions,omitempty"`
	Status               string            `json:"status,omitempty"`
	FirstMoveDone        *bool             `json:"first_move_done,omitempty"`
	RevealedSquaresCount int               `json:"revealed_squares_count,omitempty"`
}
