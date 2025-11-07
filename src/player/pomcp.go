package player

import (
	"math"
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
	"time"
)

type POMCP struct {
	root *Node
	started_playing bool
}

type Node struct {
	parent      *Node
	children    map[chess.Move]*Node
	value       float64                // v
	visits      int                    // n
	move        chess.Move             // a
	observation chess.ChessVariation   // o
	beliefs     []chess.ChessVariation // B
}

func (p *POMCP) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	time_limit := 500
	if p.root == nil {
		if board.GetNumberOfMoves() == 0 {
			p.root = &Node{nil, nil, 0, 0, chess.Move{}, nil, []chess.ChessVariation{board}}
		} else {
			var belief_state []chess.ChessVariation
			for _, move := range board.PossibleMoves() {
				follow_state := board.ExecuteMove(move)
				belief_state = append(belief_state, follow_state)
			}
			p.root = &Node{nil, nil, 0, 0, chess.Move{}, nil, belief_state}
		}
		
		p.started_playing = whiteToPlay
	}
	// TODO: prune
	selected_move := search(p.root, time_limit)
	return selected_move
}

func search(h *Node, time_limit int) chess.Move {
	start_time := time.Now()
	for !time_out(start_time, time_limit) {
		s := random_element(h.beliefs)
		simulate(s, h, 0)
	}
	return get_best_move(h)
}

func simulate(s chess.ChessVariation, h *Node, depth int) float64 {
	if math.Pow(gamma, float64(depth)) < epsilon {
		return 0
	}
	if h.children == nil { // expand
		h.children = create_all_children(s, h)
		return rollout(s, h, depth)
	}
	a := get_most_promising_action_by_ucb(s, h)
	ha := h.children[a]
	// TODO: do opponent move???

	reward := simulate(new_s, ha, depth) // TODO
	h.beliefs = append(h.beliefs, s)
	h.visits++
	ha.visits++
	ha.value = ha.value + (reward-ha.value)/float64(ha.visits)
	return reward
}

func (p POMCP) rollout(s chess.ChessVariation, h *Node, depth) float64 {
	if math.Pow(gamma, float64(depth)) < epsilon {
		return 0
	}
	new_s := state_transition(s, h)
	game_over, result := new_s.GameOver()
	if game_over {
		if !p.started_playing {
			if result == 1 {
				result = -1
			} else {
				result = 1
			}
		}
		return result // 1 for win -1 for lose
	}
	return rollout(new_s, h)
}

func random_element[T any](collection []T) T {
	selected_index := 0
	if len(collection) > 1 {
		selected_index = rand.Intn(len(collection) - 1)
	}
	return collection[selected_index]
}

func time_out(start time.Time, time_limit int) bool {
	time_out := false
	if time.Since(start).Milliseconds() > int64(time_limit) {
		time_out = true
	}
	return time_out
}

func get_best_move(h *Node) chess.Move {
	var best_node *Node
	for _, child := range h.children {
		if best_node == nil || child.value > best_node.value {
			best_node = child
		}
	}
	return best_node.move
}

func create_all_children(s chess.ChessVariation, h *Node) map[chess.Move]*Node {
	possible_transitions := make(map[chess.Move]*Node)
	// TODO: differentiate between own action and opponent action in data type
	for _, own_move := range s.PossibleMoves() {
		new_s := s.ExecuteMove(own_move)
		for _, opponent_move := range new_s.PossibleMoves() {
			new_s = new_s.ExecuteMove(opponent_move)
		}
		new_child := &Node{h, nil, 0, 0, own_move, new_s, nil}
		possible_transitions[own_move] = new_child
	}

	return possible_transitions
}

// simulates random own move and opponents
func state_transition(s chess.ChessVariation, h *Node) chess.ChessVariation {
	possible_moves := s.PossibleMoves()
	selected_move := random_element(possible_moves)
	new_s := s.ExecuteMove(selected_move)
	game_over, _ := new_s.GameOver()
	if !game_over {
		opponent_possible_moves := new_s.PossibleMoves()
		opponent_selected_move := random_element(opponent_possible_moves)
		new_s = new_s.ExecuteMove(opponent_selected_move)
	}
	return new_s
}

func get_most_promising_action_by_ucb(s chess.ChessVariation, h *Node) chess.Move {
	var max_ucb float64 = 0
	var max_child *Node
	var c float64 = 1

	for _, child := range h.children {
		ucb := child.value + c*math.Sqrt(math.Log(float64(h.visits)/float64(child.visits)))
		if max_child == nil || ucb > max_ucb {
			max_ucb = ucb
			max_child = child
		}
	}

	return max_child.move
}
