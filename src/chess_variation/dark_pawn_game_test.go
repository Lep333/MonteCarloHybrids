package chess_variation

import (
	"fmt"
	"testing"
)

func TestInitGame(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	if game.white_pawns == 0 || game.black_pawns == 0 {
		t.Errorf("Pawn Bitboards are not correctly initialized!")
	}
}

func TestGameOver(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	if game.GameOver() {
		t.Errorf("Game is in starting state and should not be over.")
	}
	game.white_pawns = 0b1 << 20
	if !game.GameOver() {
		t.Errorf("White should have won.")
	}
	game.white_pawns = 0b0
	game.black_pawns = 0b1
	if !game.GameOver() {
		t.Errorf("Black should have won.")
	}
}

func TestPossibleMovesInitialBoard(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	moves := game.PossibleMoves(true)
	if len(moves) != 5 {
		t.Errorf("Returned not 5 moves. There should be 5 possible moves...")
	}
	moves = game.PossibleMoves(false)
	if len(moves) != 5 {
		t.Errorf("Returned not 5 moves. There should be 5 possible moves...")
	}
}

func TestString(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	fmt.Printf("%v", &game)
}
