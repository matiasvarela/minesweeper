package board

import (
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper/pkg/apperrors"
	"math/rand"
	"time"
)

// NewBoard create a new board
func NewBoard(rowsNumber int, columnsNumber int, bombsNumber int) Board {
	squares := make([][]Square, rowsNumber)
	for i := range squares {
		squares[i] = make([]Square, columnsNumber)
	}

	return Board{
		Squares:        squares,
		Status:         STATUS_NEW,
		BombsPositions: &[]SquarePosition{},
		BombsNumber:    bombsNumber,
		FirstMoveDone:  newBool(false),
	}
}

// Get return the square in the given position
func (b *Board) Get(pos SquarePosition) *Square {
	return &b.Squares[pos.Row][pos.Column]
}

// Is whether square type in the given position is equals than the given type
func (b *Board) Is(pos SquarePosition, t int) bool {
	return b.Squares[pos.Row][pos.Column].Type == t
}

// GetRowsNumber return the number of rows
func (b *Board) GetRowsNumber() int {
	return len(b.Squares)
}

// GetColumnsNumber return the number of columns
func (b *Board) GetColumnsNumber() int {
	return len(b.Squares[0])
}

// GetSquaresNumber return the total number of squares
func (b *Board) GetSquaresNumber() int {
	return b.GetRowsNumber() * b.GetColumnsNumber()
}

// VerifyRange verify whether the given position is valid within the board
func (b *Board) VerifyRange(pos SquarePosition) bool {
	return pos.Row >= 0 && pos.Column >= 0 && pos.Row < b.GetRowsNumber() && pos.Column < b.GetColumnsNumber()
}

// HasNeighborBomb whether the given position has an adjacent bomb
func (b *Board) HasNeighborBomb(pos SquarePosition) bool {
	for _, n := range b.GetNeighbors(pos) {
		if b.Get(n).Type == BOMB {
			return true
		}
	}

	return false
}

// GetNeighbors return the neighbors squares for the given position
func (b *Board) GetNeighbors(pos SquarePosition) []SquarePosition {
	neighbors := []SquarePosition{}

	var target SquarePosition

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == pos.Row && j == pos.Column {
				continue
			}

			target = SquarePosition{Row: pos.Row - i, Column: pos.Column - j}

			if !b.VerifyRange(target) {
				continue
			}

			neighbors = append(neighbors, target)
		}
	}

	return neighbors
}

// FillWithBombs randomly fills the board with bombs.
// The quantity of bombs is taken from the field BombsNumber
func (b *Board) FillWithBombs(excludePosition SquarePosition) {
	rows := b.GetRowsNumber()
	columns := b.GetColumnsNumber()

	positions, extra := generateRandomPositions(b.GetSquaresNumber(), b.BombsNumber)

	var row, column int

	for _, pos := range positions {
		row = pos / columns
		column = pos - row*rows

		if row == excludePosition.Row && column == excludePosition.Column {
			row = extra / columns
			column = extra - row*rows
		}

		b.Get(SquarePosition{Row: row, Column: column}).Type = BOMB
		*b.BombsPositions = append(*b.BombsPositions, SquarePosition{Row: row, Column: column})
	}
}

// RevealSquare reveal the square in the given position trigger a reveal in cascade chain.
func (b *Board) RevealSquare(pos SquarePosition) {
	square := b.Get(pos)

	if b.Is(pos, BOMB) || b.HasNeighborBomb(pos) {
		square.Revealed = true
		b.RevealedSquaresCount++

		return
	}

	b.revealSquareInCascade(pos)
}

func (b *Board) revealSquareInCascade(pos SquarePosition) {
	if b.Is(pos, BOMB) || b.HasNeighborBomb(pos) || b.Get(pos).Revealed {
		return
	}

	b.Get(pos).Revealed = true
	b.RevealedSquaresCount++

	neighbors := b.GetNeighbors(pos)

	for i, neighbor := range neighbors {
		if b.Get(neighbor).Revealed {
			continue
		}

		b.revealSquareInCascade(neighbors[i])
	}
}

func (b *Board) PlaySquare(pos SquarePosition) error {
	if b.Status == STATUS_LOST || b.Status == STATUS_WON {
		return errors.New(apperrors.InvalidInput, nil, "cannot play a square on a finished game", "")
	}

	if !b.VerifyRange(pos) {
		return errors.New(apperrors.InvalidInput, nil, "invalid square", "")
	}

	if b.Get(pos).Revealed {
		return nil
	}

	// the first move never touch a bomb
	if !*b.FirstMoveDone {
		b.FillWithBombs(pos)
		b.Status = STATUS_ON_GOING
	}

	b.FirstMoveDone = newBool(true)

	square := b.Get(pos)
	b.RevealSquare(pos)

	switch square.Type {
	case BOMB:
		b.Status = STATUS_LOST

		for _, pos := range *b.BombsPositions {
			b.Get(pos).Marked = false
			b.Get(pos).Revealed = true
		}

		return nil
	}

	if b.RevealedSquaresCount == b.GetSquaresNumber()-b.BombsNumber {
		b.Status = STATUS_WON
	}

	return nil
}

func (b *Board) MarkSquare(pos SquarePosition) error {
	if b.Status == STATUS_LOST || b.Status == STATUS_WON {
		return errors.New(apperrors.InvalidInput, nil, "cannot mark a square on a finished game", "")
	}

	if !b.VerifyRange(pos) {
		return errors.New(apperrors.InvalidInput, nil, "invalid square", "")
	}

	if b.Get(pos).Revealed {
		return nil
	}

	b.Get(pos).Marked = !b.Get(pos).Marked

	return nil
}

// Obfuscate hide internal representation. Hide bombs positions, number of bombs, etc
func (b *Board) Obfuscate() {
	for _, pos := range *b.BombsPositions {
		if !b.Get(pos).Revealed && b.Is(pos, BOMB) {
			b.Squares[pos.Row][pos.Column].Type = EMPTY
		}
	}

	b.BombsNumber = 0
	b.BombsPositions = nil
	b.FirstMoveDone = nil
}

// ### HELPER FUNCTIONS ### //

func generateRandomPositions(n int, max int) ([]int, int) {
	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(n)

	positions := []int{}

	for _, r := range p[:max] {
		positions = append(positions, r)
	}

	return positions, p[max : max+1][0]
}

func newBool(value bool) *bool {
	return &value
}
