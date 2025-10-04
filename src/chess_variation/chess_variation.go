package chess_variation

import (
	p "monte_carlo_hybrids/player"
)

type ChessVariation interface {
	ReturnBoard(currPlayer p.Player) string
	VerifyMove(currPlayer p.Player) bool
	GameOver() bool
}
