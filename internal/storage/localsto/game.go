package localsto

import (
	"encoding/json"
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper/pkg/apperrors"

	"github.com/matiasvarela/minesweeper/internal/game"
	"git.mills.io/prologic/bitcask"
)

type GameStorage struct {
	db *bitcask.Bitcask
}

func NewGameStorage() *GameStorage {
	db, err := bitcask.Open("/tmp/minesweeper-api-db")
	if err != nil {
		panic(err)
	}

	return &GameStorage{db}
}

func (sto *GameStorage) Create(gameToCreate game.Game) error {
	bytes, err := json.Marshal(&gameToCreate)
	if err != nil {
		return errors.New(apperrors.Internal, err, "internal error", "marshal game struct into json has failed")
	}

	err = sto.db.Put([]byte(gameToCreate.ID), bytes)
	if err != nil {
		return errors.New(apperrors.Internal, err, "internal error", "put game into memory storage has failed")
	}

	return nil
}

func (sto *GameStorage) Update(gameToUpdate game.Game) error {
	has := sto.db.Has([]byte(gameToUpdate.ID))
	if !has {
		return errors.New(apperrors.NotFound, nil, "game has not been found", "game not found in memory storage")
	}

	bytes, err := json.Marshal(&gameToUpdate)
	if err != nil {
		return errors.New(apperrors.Internal, err, "internal error", "marshal game struct into json has failed")
	}

	err = sto.db.Put([]byte(gameToUpdate.ID), bytes)
	if err != nil {
		return errors.New(apperrors.Internal, err, "internal error", "save game into memory storage has failed")
	}

	return nil
}

func (sto *GameStorage) GetByID(id string) (game.Game, error) {
	has := sto.db.Has([]byte(id))
	if !has {
		return game.Game{}, errors.New(apperrors.NotFound, nil, "game has not been found", "game not found in memory storage")
	}

	bytes, err := sto.db.Get([]byte(id))
	if err != nil {
		return game.Game{}, errors.New(apperrors.Internal, err, "internal error", "get game by id from memory storage has failed")
	}

	requestedGame := game.Game{}

	err = json.Unmarshal(bytes, &requestedGame)
	if err != nil {
		return game.Game{}, errors.New(apperrors.Internal, err,"internal error", "unmarshal game into struct has failed")
	}

	return requestedGame, nil
}
