package player

import (
	"math"
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
)

type Greedy struct{}

func (g *Greedy) ViewFunc(board chess.ChessVariation) chess.ChessVariation {
	return board
}

func (g *Greedy) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	moves := board.PossibleMoves()
	best_move := chess.Move{}
	best_eval := math.Inf(-1)
	for _, move := range moves {
		new_state := board.ReturnBoard()
		new_state.ExecuteMove(move)
		eval := new_state.Heuristic(whiteToPlay)
		if over, _ := new_state.GameOver(); over {
			return move
		}
		if eval > best_eval {
			best_move = move
			best_eval = eval
		}
		if eval == best_eval && rand.Float32() >= 0.5 {
			best_move = move
			best_eval = eval
		}
	}
	return best_move
}

func (g *Greedy) String() string {
	return "greedy"
}
