package game

import (
	"crypto/rand"
	"fmt"
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper/pkg/apperrors"
	"io"
	"time"

	"github.com/matiasvarela/minesweeper/internal/board"
)

type Service interface {
	Get(id string) (Game, error)
	Create(configuration Configuration) (Game, error)
	PlaySquare(gameID string, pos board.SquarePosition) (Game, error)
	MarkSquare(gameID string, pos board.SquarePosition) (Game, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{storage}
}

func (srv *service) Get(id string) (Game, error) {
	game, err := srv.storage.GetByID(id)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return Game{}, errors.New(apperrors.NotFound, err, "game has not been found", "game not found in storage")
		}

		return Game{}, errors.New(apperrors.Internal, err, "internal error", "get game from storage has failed")
	}

	return game, nil
}

func (s *service) Create(configuration Configuration) (Game, error) {
	id, err := newUUID()
	if err != nil {
		return Game{}, errors.New(apperrors.Internal, err, "internal error", "generate new uuid has fail")
	}

	g := Game{
		ID:    id,
		Board: board.NewBoard(configuration.Rows, configuration.Columns, configuration.Bombs),
	}

	err = s.storage.Create(g)
	if err != nil {
		return Game{}, errors.New(apperrors.Internal, err, "internal error", "create game into storage has failed")
	}

	return g, nil
}

func (s *service) PlaySquare(gameID string, pos board.SquarePosition) (Game, error) {
	game, err := s.storage.GetByID(gameID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return Game{}, errors.New(apperrors.NotFound, err, "game has not been found", "game not found in storage")
		}

		return Game{}, errors.New(apperrors.Internal, err, "internal error", "get game from storage has failed")
	}

	if !*game.Board.FirstMoveDone {
		game.StartedAt = time.Now().Unix()
	}

	err = game.Board.PlaySquare(pos)
	if err != nil {
		return Game{}, errors.Wrap(err, err.Error())
	}

	err = s.storage.Update(game)
	if err != nil {
		return Game{}, errors.New(apperrors.Internal, err, "internal error", "update game into storage has failed")
	}

	return game, nil
}

func (s *service) MarkSquare(gameID string, pos board.SquarePosition) (Game, error) {
	game, err := s.storage.GetByID(gameID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return Game{}, errors.New(apperrors.NotFound, err, "game has not been found", "game not found in storage")
		}

		return Game{}, errors.New(apperrors.Internal, err, "internal error", "get game from storage has failed")
	}

	err = game.Board.MarkSquare(pos)
	if err != nil {
		return Game{}, errors.Wrap(err, err.Error())
	}

	err = s.storage.Update(game)
	if err != nil {
		return Game{}, errors.New(apperrors.Internal, err, "internal error", "update game into storage has failed")
	}

	return game, nil
}

// Helper

func newUUID() (string, error) {
	uuid := make([]byte, 16)

	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
