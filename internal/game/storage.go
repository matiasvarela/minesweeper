package game

type Storage interface {
	Create(g Game) error
	Update(g Game) error
	GetByID(id string) (Game, error)
}
