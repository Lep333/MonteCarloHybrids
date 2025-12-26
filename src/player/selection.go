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

type SelectionHybrid interface {
	Select(s chess_variation.ChessVariation) chess_variation.Move
}

type CorrectiveSelection struct {
	Bound   float64
	Epsilon float64
}

func (c *CorrectiveSelection) Select(s chess_variation.ChessVariation) chess_variation.Move {
	white := true
	if s.GetNumberOfMoves()%2 == 1 {
		white = false
	}
	default_value := s.Heuristic(white)
	moves := s.PossibleMoves()
	if len(moves) == 0 {
		return chess_variation.Move{}
	}
	score_sum := 0.0
	score := 0.0
	var move_scores []MoveScore
	for _, move := range moves {
		value := s.ExecuteMove(move).Heuristic(white)
		if value > c.Bound {
			return move
		} else if value <= default_value {
			score = c.Epsilon
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

type GreedySelection struct {
	Epsilon float64
}

// selects a move that has the best score of: heuristic or selects with x% a rand move
func (g *GreedySelection) Select(s chess_variation.ChessVariation) chess_variation.Move {
	max_score := math.Inf(-1)
	var best_move chess_variation.Move
	possible_moves := s.PossibleMoves()

	if len(possible_moves) == 0 {
		return chess_variation.Move{}
	}

	if rand.Float64() < g.Epsilon {
		return possible_moves[rand.Intn(len(possible_moves))]
	}

	white := true
	if s.GetNumberOfMoves()%2 == 1 {
		white = false
	}
	for _, move := range possible_moves {
		new_s := s.ExecuteMove(move)
		score := new_s.Heuristic(white)
		if score > max_score {
			max_score = score
			best_move = move
		}
	}
	return best_move
}

// TODO: mixed

type EarlyTermination interface {
	EarlyPlayoutTermination(s chess_variation.ChessVariation,
		depth float64) (bool, float64)
}

type EarlyPlayoutTerminationStruct struct {
	Max_depth float64
}

// returns heuristic of best move
func (e *EarlyPlayoutTerminationStruct) EarlyPlayoutTermination(
	s chess_variation.ChessVariation,
	depth float64) (bool, float64) {
	white := true
	if s.GetNumberOfMoves()%2 == 1 {
		white = false
	}

	max_score := math.Inf(-1)
	if depth < e.Max_depth {
		return false, 0.0
	} else {
		for _, move := range s.PossibleMoves() {
			new_s := s.ExecuteMove(move)
			score := new_s.Heuristic(white)
			if score > max_score {
				max_score = score
			}
		}
	}
	return true, max_score
}

// returns the next move in rollout phase
// TODO: remove determinism?
func MCTS_with_informed_rollouts(s chess_variation.ChessVariation) chess_variation.Move {
	best_move, _ := Minimax(s, true, 0, 4)
	return best_move
}

func MCTS_with_informed_cutoffs(s chess_variation.ChessVariation, depth float64, max_depth float64) (bool, float64) {
	score := math.Inf(-1)

	if depth == max_depth {
		_, score = Minimax(s, true, 0, 4)
	} else {
		return false, 0.0
	}
	return true, score
}

// inits nodes with Minimax results
func MCTS_with_informed_priors(new_node Node, s chess_variation.ChessVariation, weight int) {
	_, score := Minimax(s, true, 0, 4)
	new_node.value = score * float64(weight)
	new_node.visits = weight
}
