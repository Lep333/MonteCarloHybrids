package player

import (
	"monte_carlo_hybrids/chess_variation"
	"testing"
)

func BenchmarkSimulate(b *testing.B) {
	pomcp := POMCP{}
	dpc := chess_variation.DarkPawnChess{}
	game := dpc.ReturnBoard()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pomcp.GetMove(game, true)
	}
}
