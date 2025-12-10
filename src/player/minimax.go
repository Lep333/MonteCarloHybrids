package player

import (
	"math"
	"monte_carlo_hybrids/chess_variation"
)

func Minimax(s chess_variation.ChessVariation, max bool, depth int, depth_limit int) (chess_variation.Move, float64) {
	// TODO: check if gameover
	var best_move chess_variation.Move
	best_score := math.Inf(-1)
	game_over, _ := s.GameOver()
	if depth == depth_limit || game_over {
		empty_move := chess_variation.Move{}
		println("oink")
		return empty_move, -s.Heuristic()
	}
	for _, move := range s.PossibleMoves() {
		new_s := s.ExecuteMove(move)
		_, val := Minimax(new_s, !max, depth+1, depth_limit)
		val = -val
		if val > best_score {
			best_score = val
		}
	}
	return best_move, best_score
}
