package player

import (
	chess "monte_carlo_hybrids/chess_variation"
)

type Player interface {
	GetMove(moves []chess.Move) chess.Move
}
