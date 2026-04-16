package dto

type GameField [][]int

func NewGameField(rows, cols int) GameField {
	field := make(GameField, rows)
	for i := range field {
		field[i] = make([]int, cols)
	}
	return field
}
