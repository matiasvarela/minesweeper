package game_test

import (
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper/pkg/apperrors"
	"testing"

	"github.com/matiasvarela/minesweeper/internal/game"

	"github.com/matiasvarela/minesweeper/internal/board"
	"github.com/matiasvarela/minesweeper/internal/storage/fakesto"
	"github.com/stretchr/testify/assert"
)

var (
	service     game.Service
	fakeStorage *fakesto.GameStorage
)

func init() {
	fakeStorage = fakesto.NewGameStorage()
	service = game.NewService(fakeStorage)
}

func TestCreate(t *testing.T) {
	type input struct {
		configuration game.Configuration
	}

	tests := []struct {
		name   string
		should string
		input  input
		mock   func()
		verify func(t *testing.T, in input, g game.Game, err error)
	}{
		{
			name:   "create success",
			should: "persist the game into the storage successfully",
			input:  input{game.Configuration{Rows: 3, Columns: 3, Bombs: 2}},
			mock:   func() {},
			verify: func(t *testing.T, in input, g game.Game, err error) {
				assert.Nil(t, err)

				assert.NotEmpty(t, g.ID)
				assert.Equal(t, in.configuration.Bombs, g.Board.BombsNumber)
				assert.Equal(t, in.configuration.Rows, len(g.Board.Squares))
				assert.Equal(t, in.configuration.Columns, len(g.Board.Squares[0]))
				assert.Equal(t, board.STATUS_NEW, g.Board.Status)
			},
		},
		{
			name:   "create fails",
			should: "fail when trying to persist the game into the storage",
			input:  input{game.Configuration{Rows: 3, Columns: 3, Bombs: 2}},
			mock: func() {
				fakeStorage.AddErrorOnCreate(errors.New(apperrors.Internal, nil, "fail", ""))
			},
			verify: func(t *testing.T, in input, g game.Game, err error) {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, apperrors.Internal))
				assert.Equal(t, game.Game{}, g)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeStorage.CleanDB()
			fakeStorage.CleanErrors()

			tt.mock()

			game, err := service.Create(tt.input.configuration)

			tt.verify(t, tt.input, game, err)
		})
	}
}

func TestGet(t *testing.T) {
	type input struct {
		id string
	}

	tests := []struct {
		name   string
		should string
		input  input
		mock   func()
		verify func(t *testing.T, in input, g game.Game, err error)
	}{
		{
			name:   "get success",
			should: "get the requested game successfully",
			input:  input{"123"},
			mock: func() {
				fakeStorage.Create(game.Game{
					ID:    "123",
					Board: board.NewBoard(2, 2, 1),
				})
			},
			verify: func(t *testing.T, in input, g game.Game, err error) {
				assert.Nil(t, err)
				assert.NotEmpty(t, g.ID)
			},
		},
		{
			name:   "get fails",
			should: "fail when trying to get the requested game",
			input:  input{"123"},
			mock: func() {
				fakeStorage.AddErrorOnGetByID(errors.New(apperrors.Internal, nil, "fail", ""))
			},
			verify: func(t *testing.T, in input, g game.Game, err error) {
				assert.NotNil(t, err)
				assert.True(t, errors.Is(err, apperrors.Internal))
				assert.Equal(t, game.Game{}, g)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeStorage.CleanDB()
			fakeStorage.CleanErrors()

			tt.mock()

			game, err := service.Get(tt.input.id)

			tt.verify(t, tt.input, game, err)
		})
	}
}

func TestMarkSquare(t *testing.T) {
	type input struct {
		id  string
		pos board.SquarePosition
	}

	tests := []struct {
		name   string
		should string
		input  input
		mock   func()
		verify func(t *testing.T, in input, g game.Game, err error)
	}{
		{
			name:   "mark square success",
			should: "successfully mark a square in the board for the requested game",
			input: input{
				id:  "123",
				pos: board.SquarePosition{Row: 0, Column: 1},
			},
			mock: func() {
				fakeStorage.Create(game.Game{
					ID:    "123",
					Board: board.NewBoard(2, 2, 1),
				})
			},
			verify: func(t *testing.T, in input, g game.Game, err error) {
				assert.Nil(t, err)
				assert.NotEmpty(t, g.ID)
				assert.True(t, g.Board.Squares[0][1].Marked)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeStorage.CleanDB()
			fakeStorage.CleanErrors()

			tt.mock()

			game, err := service.MarkSquare(tt.input.id, tt.input.pos)

			tt.verify(t, tt.input, game, err)
		})
	}
}

func TestPlaySquare(t *testing.T) {
	type input struct {
		id  string
		pos board.SquarePosition
	}

	tests := []struct {
		name   string
		should string
		input  input
		mock   func()
		verify func(t *testing.T, in input, g game.Game, err error)
	}{
		{
			name:   "play square",
			should: "play a square in the requested game board",
			input: input{
				id:  "123",
				pos: board.SquarePosition{Row: 0, Column: 1},
			},
			mock: func() {
				fakeStorage.Create(game.Game{
					ID:    "123",
					Board: board.NewBoard(2, 2, 1),
				})
			},
			verify: func(t *testing.T, in input, g game.Game, err error) {
				assert.Nil(t, err)
				assert.NotEmpty(t, g.ID)
				assert.True(t, g.Board.Squares[0][1].Revealed)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeStorage.CleanDB()
			fakeStorage.CleanErrors()

			tt.mock()

			game, err := service.PlaySquare(tt.input.id, tt.input.pos)

			tt.verify(t, tt.input, game, err)
		})
	}
}
