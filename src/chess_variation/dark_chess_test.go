package chess_variation

import (
	"testing"
)

func TestDarkChessInitGame(t *testing.T) {
	game := DarkChess{}
	game.InitGame()

	if game.black_occupancy != 0xffff<<48 {
		t.Error("Black not initialized as intended.")
	}

	if game.white_occupancy != 0xffff {
		t.Error("White not initialized as intended.")
	}
}

func TestDarkChessInitPawns(t *testing.T) {
	game := DarkChess{}
	game.InitGame()
	dc_init_pawns(8)
	if dc_white_pawns_moves[8] != 16 {
		t.Error("Pawn should move forward")
	}
}
