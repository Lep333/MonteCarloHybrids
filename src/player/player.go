package player

import (
	chess "monte_carlo_hybrids/chess_variation"
)

type Player interface {
	GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move
}
