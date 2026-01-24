package player

import (
	chess "monte_carlo_hybrids/chess_variation"
)

type GreedyCaspi struct{}

func (g *GreedyCaspi) ViewFunc(board chess.ChessVariation) chess.ChessVariation {
	return board
}

func (g *GreedyCaspi) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	n := board.PossibleMoves(Moves[:])
	var captures []chess.Move
	last_pawn := Moves[0]
	for _, move := range Moves[:n] {
		if move.Capture {
			captures = append(captures, move)
		}
		if whiteToPlay && move.From < last_pawn.From {
			last_pawn = move
		}
		if !whiteToPlay && move.From > last_pawn.From {
			last_pawn = move
		}
	}
	if len(captures) > 0 {
		return random_element(captures)
	}
	return last_pawn
}

func (g *GreedyCaspi) String() string {
	return "greedy_caspi"
}
