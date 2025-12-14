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
		game.knights_moves[1] != 20992 {
		t.Errorf("Knight Bitboards are not correctly initialized!")
	}

	if game.white_king != 0b000100 || game.black_king != 0b000100<<(5*row_length_lac) {
		t.Errorf("King Bitboards are not correctly initialized!")
	}
}

func TestLosAlamosString(t *testing.T) {
	game := LosAlamosChess{}
	game.InitGame()

	expected_string := `r n q k n r 
p p p p p p 
0 0 0 0 0 0 
0 0 0 0 0 0 
P P P P P P 
R N Q K N R `
	actual_string := game.String()
	if actual_string != expected_string {
		t.Errorf("Expected String: %v but got %v", expected_string, actual_string)
	}
}

func TestLosAlamosNumberOfMoves(t *testing.T) {
	expected_moves := 3
	game := LosAlamosChess{number_of_moves: expected_moves}
	actual_moves := game.GetNumberOfMoves()

	if actual_moves != expected_moves {
		t.Errorf("Expected number of moves %v but got %v", expected_moves, actual_moves)
	}
}

func TestLosAlamosExecuteMove(t *testing.T) {
	game := LosAlamosChess{}
	game.InitGame()

	move := Move{From: 1, To: 14, Capture: false}
	new_board := game.ExecuteMove(move)
	lac, _ := new_board.(*LosAlamosChess)
	if lac.white_knights&uint(0b1<<14) == 0 {
		t.Errorf("Found no knight on field 14!")
	}

	move = Move{From: 0, To: 10, Capture: false}
	new_board = game.ExecuteMove(move)
	lac, _ = new_board.(*LosAlamosChess)
	if lac.white_rooks&uint(0b1<<10) == 0 {
		t.Errorf("Found no rook on field 10!")
	}
}

func TestLosAlamosGameOver(t *testing.T) {
	game := LosAlamosChess{}
	game.InitGame()
	game_over, result := game.GameOver()
	if game_over {
		t.Errorf("Game should not be over!")
	}
	// white wins
	game.black_king = 0
	game_over, result = game.GameOver()
	if !(game_over && result == 1) {
		t.Errorf("White should have won!")
	}

	// Black wins
	game.InitGame()
	game.white_king = 0
	game_over, result = game.GameOver()
	if !(game_over && result == -1) {
		t.Errorf("Black should have won!")
	}
}

func TestLosAlamosPossibleMoves(t *testing.T) {
	starting_move_count := 10
	game := LosAlamosChess{}
	game.InitGame()

	moves := game.PossibleMoves()
	possible_moves := len(moves)
	if possible_moves != starting_move_count {
		t.Errorf("There should be %v possible moves, got %v",
			starting_move_count, possible_moves)
	}

	move := Move{From: 6, To: 12, Capture: false}
	new_board := game.ExecuteMove(move)
	moves = new_board.PossibleMoves()
	possible_moves = len(moves)
	if possible_moves != starting_move_count {
		t.Errorf("There should be %v possible moves, got %v",
			starting_move_count, possible_moves)
	}
	move = Move{From: 24, To: 18, Capture: false}
	new_board = new_board.ExecuteMove(move)
	moves = new_board.PossibleMoves()
	possible_moves = len(moves)

	if possible_moves != starting_move_count {
		t.Errorf("There should be %v possible moves, got %v",
			9, possible_moves)
	}
}
