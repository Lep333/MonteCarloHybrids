package player

import (
	"monte_carlo_hybrids/chess_variation"
	"testing"
)

// b b b b o
// o o o o o
// o o w o o
// o w o o b
// w o o w w

// b o b o o
// o o b b o
// o o o o o
// o w o w b
// w o o o w

// b b o o o
// o o o o o
// o o o b b
// w o o o w
// o w o w o
// expect to select Move(5,11)
func TestGetMove(t *testing.T) {
	var pomcp Player
	var dark_chess chess_variation.ChessVariation
	pomcp = &POMCP{}
	dark_chess = &chess_variation.DarkPawnChess{}
	dc, ok := dark_chess.(*chess_variation.DarkPawnChess)
	if ok {
		black_pawns := 0b110000000000011 << 10
		white_pawns := 0b1000101010
		dc.InitGame()
		dc.Set_Board(uint(white_pawns), uint(black_pawns), 10, true)
		p, _ := pomcp.(*POMCP)
		p.Init_pomcp(dc, true)
		p.root = &Node{nil, nil, 0, 0, chess_variation.Move{}, nil, map[string]chess_variation.ChessVariation{}}
		p.root.beliefs[dc.String()] = dc
		p.started_playing = true
		move := p.GetMove(dark_chess, true)
		best_move := chess_variation.Move{5, 11, true}
		if move != best_move {
			t.Errorf("Best move should be %v, but returned with %v", best_move, move)
		}
	}
}
