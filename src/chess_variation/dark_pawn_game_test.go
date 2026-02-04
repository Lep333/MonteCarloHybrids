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

	game.white_pawns = 0b0
	game.black_pawns = 0b1
	game_over, winner := game.GameOver()
	if !game_over && winner == 0 {
		t.Errorf("Should be a draw.")
	}

	game.white_pawns = 0b01001
	game.black_pawns = 0b1 << 5
	game.number_of_moves = 1
	game.whiteToPlay = false
	game_over, winner = game.GameOver()
	if !game_over || winner != 1 {
		t.Errorf("Should be a win for white.")
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

func TestPossibleMovesDoublePawn(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	game.black_pawns = 0b101 << 15
	game.white_pawns = 0b100001
	moves := game.PossibleMoves()
	if len(moves) != 1 {
		t.Errorf("Returned %v moves. There should be 1 possible move...", len(moves))
	}
}

func TestPossibleMovesLeftCapture(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	game.black_pawns = 0b11111 << 20
	game.white_pawns = 0b1 << 14
	moves := game.PossibleMoves()
	if len(moves) != 1 {
		t.Errorf("Returned %v moves. There should be 1 possible move...", len(moves))
	}
}

func TestPossibleMovesRightCapture(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	game.black_pawns = 0b11111 << 20
	game.white_pawns = 0b1 << 10
	moves := game.PossibleMoves()
	if len(moves) != 1 {
		t.Errorf("Returned %v moves. There should be 1 possible move...", len(moves))
	}
}

func TestCreateView(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	view := game.CreateView(true)
	if dpc, ok := view.(*DarkPawnChess); ok {
		if dpc.black_pawns > 0 {
			t.Errorf("There should be no black pawns visible!")
		}
		fmt.Printf("%v", dpc)
		dpc.black_pawns += 0b11111 << 5
		view = dpc.CreateView(true)
		dpc, _ = view.(*DarkPawnChess)
		if dpc.black_pawns != 0b11111<<5 {
			t.Errorf("There should be 5 black pawns visible!")
		}
		fmt.Printf("%v", dpc)
		dpc.white_pawns = 0b00100
		dpc.black_pawns = 0b00100 << 5
		view = dpc.CreateView(true)
		dpc, _ = view.(*DarkPawnChess)
		if dpc.black_pawns != 0b00100<<5 {
			t.Errorf("There should be one black pawn visible!")
		}

		dpc.white_pawns = 0b100010111
		dpc.black_pawns = 0b11110 << 20
		dpc.black_pawns += 0b1 << 15
		view = dpc.CreateView(true)
		dpc, _ = view.(*DarkPawnChess)
		if dpc.black_pawns != 0 {
			t.Errorf("There should be no black pawn visible!")
		}
	}
}

func TestString(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()
	fmt.Printf("%v", &game)
}

func TestHeuristic(t *testing.T) {
	game := DarkPawnChess{}
	game.InitGame()

	expected := 0.0
	evaluation := game.Heuristic(true)
	if evaluation != expected {
		t.Errorf("Evaluation should be %v, but got %v", expected, evaluation)
	}
	game.white_pawns -= 1 << 3
	game.white_pawns += 1 << 8
	expected = 0.2
	evaluation = game.Heuristic(true)
	if evaluation != expected {
		t.Errorf("Evaluation should be %v, but got %v", expected, evaluation)
	}

	game.black_pawns -= 1 << 23
	game.black_pawns += 1 << 18
	expected = 0.0
	evaluation = game.Heuristic(false)
	if evaluation != expected {
		t.Errorf("Evaluation should be %v, but got %v", expected, evaluation)
	}
}
