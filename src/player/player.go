package player

import (
	chess "monte_carlo_hybrids/chess_variation"
)

type Player interface {
	ViewFunc(board chess.ChessVariation) chess.ChessVariation
	GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move
	String() string
}
