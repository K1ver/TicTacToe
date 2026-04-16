package model

type GameField [][]int

func NewGameField(rows, cols int) GameField {
	field := make(GameField, rows)
	for i := range field {
		field[i] = make([]int, cols)
	}
	return field
}

func (gf GameField) GetSize() (rows, cols int) {
	return len(gf), len(gf[0])
}
