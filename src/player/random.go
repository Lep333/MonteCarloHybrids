package player

import (
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
)

type RandomPlayer struct{}

func (r *RandomPlayer) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	moves := board.PossibleMoves()
	rand.Shuffle(len(moves), func(i, j int) { moves[i], moves[j] = moves[j], moves[i] })
	return moves[0]
}
