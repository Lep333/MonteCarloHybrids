package player

import (
	"math"
	"math/rand"
	"monte_carlo_hybrids/chess_variation"
	"sort"
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

type EvaluationCutOff struct {
	Threshold float64
}

func (e *EvaluationCutOff) EarlyPlayoutTermination(
	s chess_variation.ChessVariation, depth float64) (bool, float64) {
	white := true
	if s.GetNumberOfMoves()%2 == 1 {
		white = false
	}
	eval := s.Heuristic(white)
	if eval >= e.Threshold {
		return true, 1
	} else if eval <= -e.Threshold {
		return true, -1
	}
	return false, 0
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

	max_score := -1.0
	if depth < e.Max_depth {
		return false, 0.0
	} else {
		for _, move := range s.PossibleMoves() {
			new_s := s.ExecuteMove(move)
			score := new_s.Heuristic(white)
			if score >= max_score {
				max_score = score
			}
		}
	}
	return true, max_score
}

// returns the next move in rollout phase
// TODO: remove determinism?
type MCTS_with_informed_rollouts struct {
	Search_depth int
	Epsilon      float64
}

func (m *MCTS_with_informed_rollouts) Select(
	s chess_variation.ChessVariation) chess_variation.Move {
	if rand.Float64() <= m.Epsilon {
		return random_element(s.PossibleMoves())
	}
	best_move, _ := Minimax(s, true, math.Inf(-1), math.Inf(1), 0, m.Search_depth)
	return best_move
}

type MCTS_with_informed_cutoffs struct {
	Max_depth    float64
	Search_depth int
}

func (m *MCTS_with_informed_cutoffs) EarlyPlayoutTermination(
	s chess_variation.ChessVariation, depth float64) (bool, float64) {
	score := math.Inf(-1)

	if depth >= m.Max_depth {
		_, score = Minimax(s, true, math.Inf(-1), math.Inf(1), 0, m.Search_depth)
	} else {
		return false, 0.0
	}
	return true, score
}

// inits nodes with Minimax results
// Selection Phase
func MCTS_with_informed_priors(new_node Node, s chess_variation.ChessVariation, weight int) {
	_, score := Minimax(s, true, math.Inf(-1), math.Inf(1), 0, 4)
	new_node.value = score * float64(weight)
	new_node.visits = weight
}

type KBest struct {
	K int
}

func (k *KBest) Select(s chess_variation.ChessVariation) chess_variation.Move {
	moves := []MoveScore{}
	white := true
	if s.GetNumberOfMoves()%2 == 1 {
		white = false
	}
	for _, move := range s.PossibleMoves() {
		new_s := s.ExecuteMove(move)
		score := new_s.Heuristic(white)
		moves = append(moves, MoveScore{move: move, score: score})
	}
	sort.Slice(moves, func(i int, j int) bool {
		return moves[i].score > moves[j].score
	})
	if len(moves) > k.K {
		moves = moves[0:k.K]
	}
	return random_element(moves).move
}
