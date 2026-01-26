package player

import (
	"math"
	"monte_carlo_hybrids/chess_variation"
)

func Minimax(s chess_variation.ChessVariation,
	max bool, alpha float64, beta float64, depth int, depth_limit int) (chess_variation.Move, float64) {
	var best_move chess_variation.Move
	best_score := math.Inf(-1)
	game_over, _ := s.GameOver()
	white := true
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
	moves := s.PossibleMoves()
	for _, move := range moves {
		s.ExecuteMove(move)
		_, score := Minimax(s, !max, -beta, -alpha, depth+1, depth_limit)
		score = -score
		if score > best_score {
			best_score = score
			best_move = move
		}
		alpha = math.Max(alpha, score)
		if alpha >= beta {
			break
		}
	}
	return best_move, best_score
}
