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

	if game.white_king != 0b001000 || game.black_king != 0b001000<<(5*row_length_lac) {
		t.Errorf("King Bitboards are not correctly initialized!")
	}

	if game.white_queen != 0b000100 || game.black_queen != 0b000100<<(5*row_length_lac) {
		t.Errorf("Queen Bitboards are not correctly initialized!")
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

	if possible_moves != 9 {
		t.Errorf("There should be %v possible moves, got %v",
			9, possible_moves)
	}
}

func TestLosAlamosPossibleMovesRook(t *testing.T) {
	no_expected_moves := 10
	game := LosAlamosChess{white_rooks: 1, whiteToPlay: true}
	possible_moves := game.PossibleMoves()
	no_actual_moves := len(possible_moves)

	if no_actual_moves != no_expected_moves {
		t.Errorf("Expected %v moves, got %v", no_expected_moves, no_actual_moves)
	}

	no_expected_moves = 8
	game = LosAlamosChess{white_rooks: 1 << 14, black_occupancy: 1 << 20, whiteToPlay: true}
	possible_moves = game.PossibleMoves()
	no_actual_moves = len(possible_moves)

	if no_actual_moves != no_expected_moves {
		t.Errorf("Expected %v moves, got %v", no_expected_moves, no_actual_moves)
	}
}

func TestLosAlamosPossibleMovesKnight(t *testing.T) {
	no_expected_moves := 3
	game := LosAlamosChess{}
	game.InitGame()
	game.white_pawns = 0
	game.white_knights = 2
	game.white_rooks = 0
	game.white_king = 0
	game.white_queen = 0
	game.set_occupancy_boards()
	possible_moves := game.PossibleMoves()
	no_actual_moves := len(possible_moves)

	if no_actual_moves != no_expected_moves {
		t.Errorf("Expected %v moves, got %v", no_expected_moves, no_actual_moves)
	}
}

func TestLosAlamosPossibleMovesQueen(t *testing.T) {
	no_expected_moves := 17
	game := LosAlamosChess{}
	game.InitGame()
	game.white_pawns = 0
	game.white_knights = 0
	game.white_rooks = 0
	game.white_king = 0
	game.white_queen = 1 << 14
	game.set_occupancy_boards()
	possible_moves := game.PossibleMoves()
	no_actual_moves := len(possible_moves)

	if no_actual_moves != no_expected_moves {
		t.Errorf("Expected %v moves, got %v", no_expected_moves, no_actual_moves)
	}

	game.black_pawns = 1 << 19
	game.white_pawns = 1 << 8
	game.set_occupancy_boards()
	no_expected_moves = 16
	possible_moves = game.PossibleMoves()
	no_actual_moves = len(possible_moves)

	if no_actual_moves != no_expected_moves {
		t.Errorf("Expected %v moves, got %v", no_expected_moves, no_actual_moves)
	}
}

func TestLosAlamosPossibleMovesKing(t *testing.T) {
	no_expected_moves := 8
	game := LosAlamosChess{}
	game.InitGame()
	game.white_pawns = 0
	game.white_knights = 0
	game.white_rooks = 0
	game.white_king = 1 << 14
	game.white_queen = 0
	game.set_occupancy_boards()
	possible_moves := game.PossibleMoves()
	no_actual_moves := len(possible_moves)

	if no_actual_moves != no_expected_moves {
		t.Errorf("Expected %v moves, got %v", no_expected_moves, no_actual_moves)
	}

	game.black_pawns = 1 << 19
	game.white_pawns = 1 << 8
	game.set_occupancy_boards()
	no_expected_moves = 7
	possible_moves = game.PossibleMoves()
	no_actual_moves = len(possible_moves)

	if no_actual_moves != no_expected_moves {
		t.Errorf("Expected %v moves, got %v", no_expected_moves, no_actual_moves)
	}

	// r n q k r 0
	// 0 0 p p p Q
	// p 0 0 0 0 0
	// 0 0 P 0 0 0
	// P 0 0 P K P
	// R N 0 0 N R
	game = LosAlamosChess{}
	game.InitGame()

	game.white_queen += 1 << 29
	game.black_pawns -= 1 << 29
	game.set_occupancy_boards()
	move := Move{From: 1, To: 12, Capture: false}
	new_board := game.ExecuteMove(move)
	poss_moves := new_board.PossibleMoves()
	no_actual_moves = len(poss_moves)
	no_expected_moves = 10
	if no_actual_moves != no_expected_moves {
		t.Errorf("Expected %v moves, got %v", no_expected_moves, no_actual_moves)
	}
}

func TestLosAlamosCreateView(t *testing.T) {
	game := LosAlamosChess{}
	game.InitGame()

	expected_moves_no := 12
	game.white_pawns += 1 << 20
	view_board := game.CreateView()
	possible_moves := view_board.PossibleMoves()
	actual_moves_no := len(possible_moves)

	if actual_moves_no != expected_moves_no {
		t.Errorf("Expected %v moves, got %v", expected_moves_no, actual_moves_no)
	}

	expected_moves_no = 11
	game.whiteToPlay = false
	game.white_pawns -= 1 << 20
	game.black_pawns += 1 << 17
	view_board = game.CreateView()
	possible_moves = view_board.PossibleMoves()
	actual_moves_no = len(possible_moves)

	if actual_moves_no != expected_moves_no {
		t.Errorf("Expected %v moves, got %v", expected_moves_no, actual_moves_no)
	}
	// r n q k r 0
	// 0 0 p p p Q
	// p 0 0 0 0 0
	// 0 0 P 0 0 0
	// P 0 0 P K P
	// R N 0 0 N R
	game = LosAlamosChess{}
	game.InitGame()

	game.white_queen += 1 << 29
	game.black_pawns -= 1 << 29
	game.black_rooks -= 1 << 35
	game.black_rooks += 1 << 34
	game.black_knights -= 1 << 34
	game.set_occupancy_boards()
	move := Move{From: 1, To: 12, Capture: false}
	new_board := game.ExecuteMove(move).CreateView()
	poss_moves := new_board.PossibleMoves()
	no_actual_moves := len(poss_moves)
	no_expected_moves := 8
	if no_actual_moves != no_expected_moves {
		t.Errorf("Expected %v moves, got %v", no_expected_moves, no_actual_moves)
	}
}

func TestLosAlamosCreateView2(t *testing.T) {
	game := LosAlamosChess{}
	game.InitGame()

	move := Move{From: 11, To: 17, Capture: false}
	new_state := game.ExecuteMove(move)

	move = Move{From: 34, To: 23, Capture: false}
	new_state = new_state.ExecuteMove(move)

	move = Move{From: 7, To: 13, Capture: false}
	new_state = new_state.ExecuteMove(move)

	move = Move{From: 35, To: 34, Capture: false}
	new_state = new_state.ExecuteMove(move)

	move = Move{From: 10, To: 16, Capture: false}
	new_state = new_state.ExecuteMove(move)

	move = Move{From: 31, To: 18, Capture: false}
	new_state = new_state.ExecuteMove(move)

	move = Move{From: 13, To: 18, Capture: true}
	new_state = new_state.ExecuteMove(move)

	move = Move{From: 23, To: 10, Capture: false}
	new_state = new_state.ExecuteMove(move)

	move = Move{From: 3, To: 10, Capture: true}
	new_state = new_state.ExecuteMove(move)
	expected_possible_moves := len(new_state.PossibleMoves())
	moves := new_state.CreateView().PossibleMoves()
	actual_possible_moves := len(moves)
	if expected_possible_moves != actual_possible_moves {
		t.Errorf("Error")
	}
}
