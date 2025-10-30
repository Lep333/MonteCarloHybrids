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
	game_over, _ := game.GameOver()
	if game_over {
		t.Errorf("Game is in starting state and should not be over.")
	}
	game.white_pawns = 0b1 << 20
	game_over, _ = game.GameOver()
	if !game_over {
		t.Errorf("White should have won.")
	}
	game.white_pawns = 0b0
	game.black_pawns = 0b1
	game_over, _ = game.GameOver()
	if !game_over {
		t.Errorf("Black should have won.")
	}
}

func TestPossibleMovesInitialBoard(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	moves := game.PossibleMoves()
	if len(moves) != 5 {
		t.Errorf("Returned not 5 moves. There should be 5 possible moves...")
	}
	moves = game.PossibleMoves()
	if len(moves) != 5 {
		t.Errorf("Returned not 5 moves. There should be 5 possible moves...")
	}
}

func TestPossibleMoves(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	game.black_pawns = 0b101 << 5
	game.white_pawns = 0b10
	moves := game.PossibleMoves()
	if len(moves) != 3 {
		t.Errorf("Returned %v moves. There should be 3 possible moves...", len(moves))
	}
}

func TestCreateView(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	view := game.CreateView()
	if dpc, ok := view.(*DarkPawnChess); ok {
		if dpc.black_pawns > 0 {
			t.Errorf("There should be no black pawns visible!")
		}
		fmt.Printf("%v", dpc)
		dpc.black_pawns += 0b11111 << 5
		view = dpc.CreateView()
		dpc, _ = view.(*DarkPawnChess)
		if dpc.black_pawns != 0b11111<<5 {
			t.Errorf("There should be 5 black pawns visible!")
		}
		fmt.Printf("%v", dpc)
	}

}

func TestString(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	fmt.Printf("%v", &game)
}
