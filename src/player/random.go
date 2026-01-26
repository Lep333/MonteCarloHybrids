package player

import (
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
)

type RandomPlayer struct{}

func (r *RandomPlayer) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	moves := board.PossibleMoves()
	for _, move := range board.PossibleMoves() {
		if move.Capture {
			return move
		}
	}
	return moves[rand.Intn(len(moves))]
}

func (r *RandomPlayer) String() string {
	return "RandomPlayer"
}
