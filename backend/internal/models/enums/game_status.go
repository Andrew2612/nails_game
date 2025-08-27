package enums

type GameStatus int

const (
	InProgress GameStatus = iota
	Draw
	FirstPlayerWon
	SecondPlayerWon
)

func (s GameStatus) String() string {
	return [...]string{"CREATED", "IN_PROGRESS", "FINISHED"}[s]
}
