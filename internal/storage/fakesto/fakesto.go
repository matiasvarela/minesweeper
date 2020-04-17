package fakesto

import (
	"encoding/json"
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper/pkg/apperrors"

	"github.com/matiasvarela/minesweeper/internal/game"
)

type GameStorage struct {
	db     map[string][]byte
	errors map[string]error
}

func NewGameStorage() *GameStorage {
	return &GameStorage{
		map[string][]byte{},
		map[string]error{},
	}
}

func (sto *GameStorage) CleanErrors() {
	sto.errors = map[string]error{}
}

func (sto *GameStorage) CleanDB() {
	sto.db = map[string][]byte{}
}

func (sto *GameStorage) AddErrorOnCreate(err error) {
	sto.errors["on_create"] = err
}

func (sto *GameStorage) AddErrorOnUpdate(err error) {
	sto.errors["on_update"] = err
}

func (sto *GameStorage) AddErrorOnGetByID(err error) {
	sto.errors["on_get_by_id"] = err
}

func (sto *GameStorage) Create(gameToCreate game.Game) error {
	if err, ok := sto.errors["on_create"]; ok {
		return err
	}

	bytes, err := json.Marshal(&gameToCreate)
	if err != nil {
		return errors.New(apperrors.Internal,err,"internal error", "marshal game struct into json has failed")
	}

	sto.db[gameToCreate.ID] = bytes

	return nil
}

func (sto *GameStorage) Update(gameToUpdate game.Game) error {
	if err, ok := sto.errors["on_update"]; ok {
		return err
	}

	if _, ok := sto.db[gameToUpdate.ID]; !ok {
		return errors.New(apperrors.NotFound, nil, "game has not been found", "game not found in db")
	}

	bytes, err := json.Marshal(&gameToUpdate)
	if err != nil {
		return errors.New(apperrors.Internal, err, "internal error", "marshal game struct into json has failed")
	}

	sto.db[gameToUpdate.ID] = bytes

	return nil
}

func (sto *GameStorage) GetByID(id string) (game.Game, error) {
	if err, ok := sto.errors["on_get_by_id"]; ok {
		return game.Game{}, err
	}

	if _, ok := sto.db[id]; !ok {
		return game.Game{}, errors.New(apperrors.NotFound, nil, "game has not been found", "game not found in db")
	}

	bytes := sto.db[id]

	requestedGame := game.Game{}

	err := json.Unmarshal(bytes, &requestedGame)
	if err != nil {
		return game.Game{}, errors.New(apperrors.Internal, err,"internal error", "unmarshal game into struct has failed")
	}

	return requestedGame, nil
}