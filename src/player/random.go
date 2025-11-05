package player

import (
	chess "monte_carlo_hybrids/chess_variation"
	"sort"
)

type RandomPlayer struct{}

func (r *RandomPlayer) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	moves := board.PossibleMoves()
	sort.Slice(moves, func(i, j int) bool {
		return moves[i].Capture && !moves[j].Capture
	})
	// rand.Shuffle(len(moves), func(i, j int) { moves[i], moves[j] = moves[j], moves[i] })
	return moves[0]
}
