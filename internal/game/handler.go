package game

import (
	"github.com/gin-gonic/gin"
	"github.com/matiasvarela/errors"
	"github.com/matiasvarela/minesweeper/internal/board"
	"github.com/matiasvarela/minesweeper/pkg/apperrors"
	"gopkg.in/go-playground/validator.v8"
)

type HttpHandler interface {
	Create(*gin.Context)
	Get(*gin.Context)
	PlaySquare(c *gin.Context)
	MarkSquare(c *gin.Context)
}

type httpHandler struct {
	service Service
}

var (
	validate *validator.Validate
)

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
	validate.RegisterStructValidation(ConfigurationStructValidation, Configuration{})
}

func NewHttpHandler(service Service) HttpHandler {
	return &httpHandler{service}
}

func (h *httpHandler) Get(c *gin.Context) {
	game, err := h.service.Get(c.Param("id"))
	if err != nil {
		apierr := apperrors.ToApiError(err)
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	game.updateElapsedTime()
	game.Board.Obfuscate()

	c.JSON(200, game)
}

func (h *httpHandler) Create(c *gin.Context) {
	configuration := Configuration{}

	err := c.BindJSON(&configuration)
	if err != nil {
		apierr := apperrors.ToApiError(errors.New(apperrors.InvalidInput, err, "invalid body", "bind json has failed"))
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	err = validate.Struct(configuration)
	if err != nil {
		apierr := apperrors.ToApiError(errors.New(apperrors.InvalidInput, err, "invalid body", "validations has failed"))
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	game, err := h.service.Create(configuration)
	if err != nil {
		apierr := apperrors.ToApiError(err)
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	game.updateElapsedTime()
	game.Board.Obfuscate()

	c.JSON(201, game)
}

func (h *httpHandler) PlaySquare(c *gin.Context) {
	body := PlaySquareBody{}

	err := c.BindJSON(&body)
	if err != nil {
		apierr := apperrors.ToApiError(errors.New(apperrors.InvalidInput, err, "invalid body", "bind json has failed"))
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	err = validate.Struct(body)
	if err != nil {
		apierr := apperrors.ToApiError(errors.New(apperrors.InvalidInput, err, "invalid body", "validations has failed"))
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	game, err := h.service.PlaySquare(c.Param("id"), board.SquarePosition{Row: body.Row, Column: body.Column})
	if err != nil {
		apierr := apperrors.ToApiError(err)
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	game.updateElapsedTime()
	game.Board.Obfuscate()

	c.JSON(200, game)
}

func (h *httpHandler) MarkSquare(c *gin.Context) {
	body := MarkSquareBody{}

	err := c.BindJSON(&body)
	if err != nil {
		apierr := apperrors.ToApiError(errors.New(apperrors.InvalidInput, err, "invalid body", "bind json has failed"))
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	err = validate.Struct(body)
	if err != nil {
		apierr := apperrors.ToApiError(errors.New(apperrors.InvalidInput, err, "invalid body", "validations has failed"))
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	game, err := h.service.MarkSquare(c.Param("id"), board.SquarePosition{Row: body.Row, Column: body.Column})
	if err != nil {
		apierr := apperrors.ToApiError(err)
		c.AbortWithStatusJSON(apierr.Status, apierr)
		return
	}

	game.updateElapsedTime()
	game.Board.Obfuscate()

	c.JSON(200, game)
}