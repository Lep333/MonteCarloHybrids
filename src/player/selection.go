package player

import (
	"math"
	"math/rand"
	"monte_carlo_hybrids/chess_variation"
)

type MoveScore struct {
	move  chess_variation.Move
	score float64
}

func Corrective(s chess_variation.ChessVariation, h *Node) chess_variation.Move {
	default_value := s.Heuristic()
	moves := s.PossibleMoves()
	if len(moves) == 0 {
		return chess_variation.Move{}
	}
	score_sum := 0.0
	bound := 0.95 // should be needed to be a win
	epsilon := 0.0095
	score := 0.0
	var move_scores []MoveScore
	for _, move := range moves {
		value := s.ExecuteMove(move).Heuristic()
		if value > bound {
			return move
		} else if value <= default_value {
			score = epsilon
		} else {
			score = value
		}
		move_scores = append(move_scores, MoveScore{move, score})
		score_sum += score
	}
	score_sum *= rand.Float64()
	for _, move := range move_scores {
		score_sum -= move.score
		if score_sum < 0 {
			return move.move
		}
	}
	return moves[len(moves)-1]
}

// selects a move that has the best score of: heuristic or selects with x% a rand move
func Greedy(s chess_variation.ChessVariation, h *Node) chess_variation.Move {
	rand_move_probability := 0.2
	max_score := math.Inf(-1)
	var best_move chess_variation.Move
	possible_moves := s.PossibleMoves()

	if len(possible_moves) == 0 {
		return chess_variation.Move{}
	}

	if rand.Float64() < rand_move_probability {
		return possible_moves[rand.Intn(len(possible_moves))]
	}

	for _, move := range possible_moves {
		new_s := s.ExecuteMove(move)
		score := new_s.Heuristic()
		if score > max_score {
			max_score = score
			best_move = move
		}
	}
	return best_move
}
