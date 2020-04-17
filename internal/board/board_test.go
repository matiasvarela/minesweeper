package board_test

import (
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper/pkg/apperrors"
	"testing"

	"github.com/matiasvarela/minesweeper/internal/board"
	"github.com/stretchr/testify/assert"
)

var (
	mocks boardMocks
)

type boardMocks map[string]board.Board

func (m boardMocks) get(id string) *board.Board {
	board := m[id]

	return &board
}

func init() {
	_true := true
	_false := false

	mocks = boardMocks{
		"lost_board": board.Board{
			Status: board.STATUS_LOST,
			Squares: [][]board.Square{
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: true}, {Type: board.EMPTY, Revealed: false}},
				{{Type: board.EMPTY, Revealed: true}, {Type: board.BOMB, Revealed: true}, {Type: board.EMPTY, Revealed: true}},
				{{Type: board.EMPTY, Revealed: true}, {Type: board.BOMB, Revealed: true}, {Type: board.EMPTY, Revealed: true}},
			},
			BombsPositions:       &[]board.SquarePosition{{1, 1}, {2, 1}},
			FirstMoveDone:        &_true,
			RevealedSquaresCount: 6,
			BombsNumber: 2,
		},
		"won_board": board.Board{
			Status: board.STATUS_WON,
			Squares: [][]board.Square{
				{{Type: board.EMPTY, Revealed: true}, {Type: board.EMPTY, Revealed: true}, {Type: board.EMPTY, Revealed: true}},
				{{Type: board.EMPTY, Revealed: true}, {Type: board.BOMB, Revealed: false}, {Type: board.EMPTY, Revealed: true}},
				{{Type: board.EMPTY, Revealed: true}, {Type: board.BOMB, Revealed: false}, {Type: board.EMPTY, Revealed: true}},
			},
			BombsPositions:       &[]board.SquarePosition{{1, 1}, {2, 1}},
			FirstMoveDone:        &_true,
			RevealedSquaresCount: 7,
			BombsNumber: 2,
		},
		"on_going_board": board.Board{
			Status: board.STATUS_ON_GOING,
			Squares: [][]board.Square{
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: true}, {Type: board.EMPTY, Revealed: false}},
				{{Type: board.EMPTY, Revealed: true}, {Type: board.BOMB, Revealed: false}, {Type: board.EMPTY, Revealed: true}},
				{{Type: board.EMPTY, Revealed: false, Marked: true}, {Type: board.BOMB, Revealed: false}, {Type: board.EMPTY, Revealed: true}},
			},
			BombsPositions:       &[]board.SquarePosition{{1, 1}, {2, 1}},
			FirstMoveDone:        &_true,
			RevealedSquaresCount: 4,
			BombsNumber: 2,
		},
		"last_to_win_board": board.Board{
			Status: board.STATUS_ON_GOING,
			Squares: [][]board.Square{
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: true}, {Type: board.EMPTY, Revealed: true}},
				{{Type: board.EMPTY, Revealed: true}, {Type: board.BOMB, Revealed: false}, {Type: board.EMPTY, Revealed: true}},
				{{Type: board.EMPTY, Revealed: true}, {Type: board.BOMB, Revealed: false}, {Type: board.EMPTY, Revealed: true}},
			},
			BombsPositions:       &[]board.SquarePosition{{1, 1}, {2, 1}},
			FirstMoveDone:        &_true,
			RevealedSquaresCount: 6,
			BombsNumber: 2,
		},
		"reveal_in_cascade_board": board.Board{
			Status: board.STATUS_ON_GOING,
			Squares: [][]board.Square{
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}},
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}},
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}},
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}, {Type: board.BOMB, Revealed: false}},
			},
			BombsPositions:       &[]board.SquarePosition{{3, 2}},
			FirstMoveDone:        &_true,
			RevealedSquaresCount: 0,
			BombsNumber:          1,
		},
		"new_board": board.Board{
			Status: board.STATUS_NEW,
			Squares: [][]board.Square{
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}},
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}},
				{{Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}, {Type: board.EMPTY, Revealed: false}},
			},
			BombsPositions:       &[]board.SquarePosition{},
			FirstMoveDone:        &_false,
			RevealedSquaresCount: 0,
			BombsNumber: 3,
		},
	}
}

func TestBoard_MarkSquare(t *testing.T) {
	type input struct {
		board *board.Board
		pos   board.SquarePosition
	}

	tests := []struct {
		name   string
		should string
		input  input
		verify func(t *testing.T, in input, err error)
	}{
		{
			name:   "mark successfully",
			should: "mark the request square",
			input: input{
				board: mocks.get("on_going_board"),
				pos:   board.SquarePosition{0, 0},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.Nil(t, err)
				assert.True(t, in.board.Squares[in.pos.Row][in.pos.Column].Marked)
			},
		},
		{
			name:   "unmark successfully",
			should: "unmark the request square",
			input: input{
				board: mocks.get("on_going_board"),
				pos:   board.SquarePosition{2, 0},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.Nil(t, err)
				assert.False(t, in.board.Squares[in.pos.Row][in.pos.Column].Marked)
			},
		},
		{
			name:   "out of range",
			should: "return an invalid input error",
			input: input{
				board: mocks.get("on_going_board"),
				pos:   board.SquarePosition{10, 10},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, apperrors.InvalidInput))
			},
		},
		{
			name:   "mark on a revealed square",
			should: "do nothing",
			input: input{
				board: mocks.get("on_going_board"),
				pos:   board.SquarePosition{0, 1},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.Nil(t, err)
				assert.False(t, in.board.Squares[in.pos.Row][in.pos.Column].Marked)
			},
		},
		{
			name:   "test mark on a finished board",
			should: "return an invalid input error",
			input: input{
				board: mocks.get("lost_board"),
				pos:   board.SquarePosition{0, 0},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, apperrors.InvalidInput))
			},
		},
		{
			name:   "test mark on a finished board",
			should: "return an invalid input error",
			input: input{
				board: mocks.get("won_board"),
				pos:   board.SquarePosition{0, 0},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, apperrors.InvalidInput))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.board.MarkSquare(tt.input.pos)
			tt.verify(t, tt.input, err)
		})
	}
}

func TestBoard_PlaySquare(t *testing.T) {
	type input struct {
		board *board.Board
		pos   board.SquarePosition
	}

	tests := []struct {
		name   string
		should string
		input  input
		verify func(t *testing.T, in input, err error)
	}{
		{
			name:   "reveal square successfully",
			should: "reveal the request square",
			input: input{
				board: mocks.get("on_going_board"),
				pos:   board.SquarePosition{0, 0},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.Nil(t, err)
				assert.True(t, in.board.Squares[in.pos.Row][in.pos.Column].Revealed)
				assert.Equal(t, 5, in.board.RevealedSquaresCount)
			},
		},
		{
			name:   "reveal square with bomb successfully",
			should: "lost the game",
			input: input{
				board: mocks.get("on_going_board"),
				pos:   board.SquarePosition{1, 1},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.Nil(t, err)
				assert.True(t, in.board.Squares[in.pos.Row][in.pos.Column].Revealed)
				assert.Equal(t, board.STATUS_LOST, in.board.Status)
			},
		},
		{
			name:   "reveal last square successfully",
			should: "win the game",
			input: input{
				board: mocks.get("last_to_win_board"),
				pos:   board.SquarePosition{0, 0},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.Nil(t, err)
				assert.True(t, in.board.Squares[in.pos.Row][in.pos.Column].Revealed)
				assert.Equal(t, board.STATUS_WON, in.board.Status)
			},
		},
		{
			name:   "reveal square in cascade successfully",
			should: "reveal many squares",
			input: input{
				board: mocks.get("reveal_in_cascade_board"),
				pos:   board.SquarePosition{0, 0},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.Nil(t, err)
				assert.True(t, in.board.Squares[in.pos.Row][in.pos.Column].Revealed)
				assert.Equal(t, board.STATUS_ON_GOING, in.board.Status)
				assert.Equal(t, 8, in.board.RevealedSquaresCount)
			},
		},
		{
			name:   "first play",
			should: "fill with bombs and reveal square",
			input: input{
				board: mocks.get("new_board"),
				pos:   board.SquarePosition{0, 0},
			},
			verify: func(t *testing.T, in input, err error) {
				assert.Nil(t, err)
				assert.True(t, in.board.Squares[in.pos.Row][in.pos.Column].Revealed)
				assert.Equal(t, board.STATUS_ON_GOING, in.board.Status)
				assert.True(t, len(*in.board.BombsPositions) > 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.board.PlaySquare(tt.input.pos)
			tt.verify(t, tt.input, err)
		})
	}
}
