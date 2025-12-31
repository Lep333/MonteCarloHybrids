package player

import (
	"math"
	"monte_carlo_hybrids/chess_variation"
)

func Minimax(s chess_variation.ChessVariation,
	max bool, alpha float64, beta float64, depth int, depth_limit int) (chess_variation.Move, float64) {
	// TODO: check if gameover
	var best_move chess_variation.Move
	best_score := math.Inf(-1)
	game_over, _ := s.GameOver()
	white := true
	value := math.Inf(-1)
	if s.GetNumberOfMoves()%2 == 1 {
		white = false
	}
	if depth == depth_limit {
		empty_move := chess_variation.Move{}
		return empty_move, s.Heuristic(!white)
	}
	if game_over {
		empty_move := chess_variation.Move{}
		return empty_move, s.Heuristic(white)
	}
	for _, move := range s.PossibleMoves() {
		new_s := s.ExecuteMove(move)
		_, val := Minimax(new_s, !max, -beta, -alpha, depth+1, depth_limit)
		value = math.Max(val, value)
		val = -val
		alpha = math.Max(alpha, value)
		if alpha >= beta {
			break
		}
		if val > best_score {
			best_score = val
		}
	}
	return best_move, best_score
}
