package player

import (
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
)

type RandomPlayer struct{}

func (r *RandomPlayer) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	n := board.PossibleMoves(Moves[:])
	for _, move := range Moves[:n] {
		if move.Capture {
			return move
		}
	}
	return Moves[rand.Intn(n)]
}

func (r *RandomPlayer) String() string {
	return "RandomPlayer"
}
