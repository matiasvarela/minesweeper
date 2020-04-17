package game

import (
	"reflect"
	"time"

	"github.com/matiasvarela/minesweeper/internal/board"
	"gopkg.in/go-playground/validator.v8"
)

type Game struct {
	ID          string      `json:"id"`
	Board       board.Board `json:"board"`
	StartedAt   int64       `json:"started_at"`
	ElapsedTime int64       `json:"elapsed_time"`
}

func (g *Game) updateElapsedTime() {
	if g.StartedAt > 0 {
		g.ElapsedTime = time.Now().Unix() - g.StartedAt
	}
}

type Configuration struct {
	Rows    int `json:"rows" validate:"required,gte=3"`
	Columns int `json:"columns" validate:"required,gte=3"`
	Bombs   int `json:"bombs" validate:"required,gte=0"`
}

type PlaySquareBody struct {
	Row    int `json:"row" validate:"gte=0"`
	Column int `json:"column" validate:"gte=0"`
}

type MarkSquareBody struct {
	Row    int `json:"row" validate:"gte=0"`
	Column int `json:"column" validate:"gte=0"`
}

func ConfigurationStructValidation(v *validator.Validate, structLevel *validator.StructLevel) {
	configuration := structLevel.CurrentStruct.Interface().(Configuration)

	if configuration.Bombs >= (configuration.Rows*configuration.Columns)-1 {
		structLevel.ReportError(reflect.ValueOf(configuration.Bombs), "Bombs", "bombs", "bombsnumber")
	}
}
