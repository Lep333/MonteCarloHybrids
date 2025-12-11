package chess_variation

import (
	"testing"
)

func TestLosAlamosInitGame(t *testing.T) {
	game := LosAlamosChess{}
	game.InitGame()

	if game.white_pawns == 0 || game.black_pawns == 0 {
		t.Errorf("Pawn Bitboards are not correctly initialized!")
	}

	if game.white_rooks != 0b100001 || game.black_rooks != 0b100001<<(5*row_length_lac) {
		t.Errorf("Rook Bitboards are not correctly initialized!")
	}

	if game.white_knights != 0b010010 ||
		game.black_knights != 0b010010<<(5*row_length_lac) ||
		game.knights_moves[0] != 21 {
		t.Errorf("Knight Bitboards are not correctly initialized!")
	}

	if game.white_king != 0b000100 || game.black_king != 0b000100<<(5*row_length_lac) {
		t.Errorf("King Bitboards are not correctly initialized!")
	}
}
