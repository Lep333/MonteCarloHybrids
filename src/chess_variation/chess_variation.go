package chess_variation

import (
	p "monte_carlo_hybrids/player"
)

type Move struct {
	from int8
	to   int8
}

type ChessVariation interface {
	InitGame()
	ReturnBoard(currPlayer p.Player) string
	VerifyMove(currPlayer p.Player) bool
	PossibleMoves(player int) []Move
	GameOver() bool
}
