package player

import (
	"math"
	"monte_carlo_hybrids/chess_variation"
)

type Move_With_Score struct {
	move  chess_variation.Move
	score float64
}

func Minimax(s chess_variation.ChessVariation, mini bool, depth int) []Move_With_Score {
	var moves_with_scores []Move_With_Score
	if depth == 10 {
		return []Move_With_Score{{chess_variation.Move{}, 1}}
	}
	possible_moves := s.PossibleMoves()
	for _, move := range possible_moves {
		new_s := s.ExecuteMove(move)
		game_over, _ := new_s.GameOver()
		if game_over {
			result := 0
			if mini {
				result = -1
			} else {
				result = 1
			}
			return []Move_With_Score{{move, float64(result)}}
		}
		move_scores := Minimax(new_s, !mini, depth+1)
		var best_move Move_With_Score
		if mini {
			lowest := math.Inf(1)
			for _, scores := range move_scores {
				if scores.score < lowest {
					lowest = scores.score
					best_move = scores
				}
			}
		} else {
			highest := math.Inf(-1)
			for _, scores := range move_scores {
				if scores.score > highest {
					highest = scores.score
					best_move = scores
				}
			}
		}
		best_move.move = move
		moves_with_scores = append(moves_with_scores, best_move)
	}
	return moves_with_scores
}
