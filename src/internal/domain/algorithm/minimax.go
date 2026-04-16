package algorithm

import (
	"app/internal/domain/model"
)

type MinimaxAlgorithm interface {
	FindBestMove(field model.GameField, player int, maxDepth int) (int, int, error)
}

type minimax struct {
}

func NewMinimax() MinimaxAlgorithm {
	return &minimax{}
}

func (m *minimax) FindBestMove(field model.GameField, player int, maxDepth int) (int, int, error) {
	bestScore := -1 << 31
	bestRow, bestCol := -1, -1

	rows, cols := field.GetSize()

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if field[i][j] == 0 {
				field[i][j] = player
				score := m.minimax(field, 0, false, player, maxDepth)
				field[i][j] = 0

				if score > bestScore {
					bestScore = score
					bestRow, bestCol = i, j
				}
			}
		}
	}

	return bestRow, bestCol, nil
}

func (m *minimax) minimax(field model.GameField, depth int, isMax bool, player int, maxDepth int) int {
	if depth == maxDepth || m.isTerminal(field) {
		return m.evaluate(field, player)
	}

	opponent := 3 - player
	bestScore := -1 << 31
	if !isMax {
		bestScore = 1 << 31
	}

	rows, cols := field.GetSize()

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if field[i][j] == 0 {
				currentPlayer := player
				if !isMax {
					currentPlayer = opponent
				}

				field[i][j] = currentPlayer
				score := m.minimax(field, depth+1, !isMax, player, maxDepth)
				field[i][j] = 0

				if isMax {
					if score > bestScore {
						bestScore = score
					}
				} else {
					if score < bestScore {
						bestScore = score
					}
				}
			}
		}
	}

	return bestScore
}

func (m *minimax) evaluate(field model.GameField, player int) int {
	// Упрощенная функция оценки
	opponent := 3 - player
	score := 0

	// Проверяем выигрышные комбинации
	if m.checkWinner(field, player) {
		return 1000
	}
	if m.checkWinner(field, opponent) {
		return -1000
	}

	return score
}

func (m *minimax) isTerminal(field model.GameField) bool {
	return m.checkWinner(field, 1) || m.checkWinner(field, 2) || m.isBoardFull(field)
}

func (m *minimax) checkWinner(field model.GameField, player int) bool {
	rows, cols := field.GetSize()

	// Проверка строк и столбцов
	for i := 0; i < rows; i++ {
		for j := 0; j <= cols-3; j++ {
			if field[i][j] == player && field[i][j+1] == player && field[i][j+2] == player {
				return true
			}
		}
	}

	for j := 0; j < cols; j++ {
		for i := 0; i <= rows-3; i++ {
			if field[i][j] == player && field[i+1][j] == player && field[i+2][j] == player {
				return true
			}
		}
	}

	// Проверка главных диагоналей (сверху-слева направо-вниз)
	for i := 0; i <= rows-3; i++ {
		for j := 0; j <= cols-3; j++ {
			if field[i][j] == player && field[i+1][j+1] == player && field[i+2][j+2] == player {
				return true
			}
		}
	}

	// Проверка побочных диагоналей (сверху-справа налево-вниз)
	for i := 0; i <= rows-3; i++ {
		for j := 2; j < cols; j++ {
			if field[i][j] == player && field[i+1][j-1] == player && field[i+2][j-2] == player {
				return true
			}
		}
	}

	return false
}

func (m *minimax) isBoardFull(field model.GameField) bool {
	rows, cols := field.GetSize()
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if field[i][j] == 0 {
				return false
			}
		}
	}
	return true
}
