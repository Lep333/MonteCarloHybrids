package player

import (
	"math"
	"monte_carlo_hybrids/chess_variation"
	"testing"
)

// b b o o o
// o o o o o
// o o o b b
// w o o o w
// o w o w o
// expect to select Move(5,11)
func TestGetMoveMinimax(t *testing.T) {
	var pomcp Player
	var dark_chess chess_variation.ChessVariation
	settings := Settings{
		Termination_parameter: 10000,
		Gamma:                 0.95,
		Epsilon:               0.005,
		Ucb_c:                 1,
		Capture_reward:        0.2,
		Rollout_capture:       0.2,
	}
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
		p.Root = &Node{nil, nil, 0, 0, chess_variation.Move{}, nil, map[string]chess_variation.ChessVariation{}}
		p.Root.beliefs[dc.String()] = dc
		p.Started_playing = true
		p.Settings = settings
		moves := Minimax(dark_chess, false, 0)
		best_move := chess_variation.Move{5, 11, true}
		minimax_best_move := Move_With_Score{chess_variation.Move{}, math.Inf(-1)}
		for _, move := range moves {
			if move.score > minimax_best_move.score {
				minimax_best_move = move
			}
		}
		if minimax_best_move.move != best_move {
			t.Errorf("Best move should be %v, but returned with %v", best_move, minimax_best_move.move)
		}
	}
}
