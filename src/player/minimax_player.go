package player

import (
	"math"
	chess "monte_carlo_hybrids/chess_variation"
)

type Mini struct{}

func (g *Mini) ViewFunc(board chess.ChessVariation) chess.ChessVariation {
	return board
}

func (g *Mini) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	best_move, _ := Minimax(board, true, math.Inf(-1), math.Inf(1), 0, 2)
	return best_move
}

func (g *Mini) String() string {
	return "minimax_player_depth_2"
}
