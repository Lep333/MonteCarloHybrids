package player

import (
	"math"
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
	"time"
)

const gamma = 0.95
const epsilon = 0.005

type POMCP struct {
	root            *Node
	started_playing bool
	last_move       chess.Move
}

type Combined_Key struct {
	action      chess.Move
	observation string
}

type Node struct {
	parent      *Node
	children    map[Combined_Key]*Node
	value       float64                // v
	visits      int                    // n
	move        chess.Move             // a
	observation chess.ChessVariation   // o
	beliefs     []chess.ChessVariation // B
}

func (p *POMCP) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	time_limit := 1000

	if p.root == nil {
		if board.GetNumberOfMoves() == 0 {
			p.root = &Node{nil, nil, 0, 0, chess.Move{}, nil, []chess.ChessVariation{board}}
		} else if board.GetNumberOfMoves() == 1 {
			var belief_state []chess.ChessVariation
			init_board := board.ReturnBoard()
			init_board.InitGame()
			for _, move := range init_board.PossibleMoves() { // TODO: board with state 0
				follow_state := init_board.ExecuteMove(move)
				belief_state = append(belief_state, follow_state)
			}
			p.root = &Node{nil, nil, 0, 0, chess.Move{}, nil, belief_state}
		}
		p.started_playing = whiteToPlay
	}
	// prune
	if p.last_move != (chess.Move{}) {
		board_string := board.String()
		new_root := p.root.children[Combined_Key{p.last_move, board_string}]
		if new_root == nil {
			// TODO: init root with consistent beliefs
			println("ohoh")
			for key, _ := range p.root.children {
				println(key.action.From, key.action.To)
				print(key.observation)
				println("______")
			}
		}
		p.root = new_root // TODO: find
		// update consistent beliefs
		consistent_beliefs := []chess.ChessVariation{}
		for _, belief := range p.root.beliefs {
			if belief.CreateView().String() == board.String() {
				consistent_beliefs = append(consistent_beliefs, belief)
			}
		}
		p.root.beliefs = consistent_beliefs
	}
	for _, belief := range p.root.beliefs { // TODO: Root with all action observation pairs init
		for _, action := range belief.PossibleMoves() {
			// TODO: enemy move
			if p.root.children == nil {
				p.root.children = make(map[Combined_Key]*Node)
			}
			new_board := belief.ExecuteMove(action)
			for _, move := range new_board.PossibleMoves() { // TODO what if no possible move?
				opp_move := new_board.ExecuteMove(move)
				elem := p.root.children[Combined_Key{action, opp_move.CreateView().String()}]
				if elem == nil {
					p.root.children[Combined_Key{action, opp_move.CreateView().String()}] = &Node{p.root, nil, 0, 0, action, opp_move, nil}
				}
			}
		}
	}
	selected_move := search(p.root, time_limit, p)
	println("visits: ", p.root.visits)
	p.last_move = selected_move
	return selected_move
}

func search(h *Node, time_limit int, p *POMCP) chess.Move {
	start_time := time.Now()
	for !time_out(start_time, time_limit) {
		s := random_element(h.beliefs)
		simulate(s, h, 0, p)
	}
	return get_best_move(h)
}

func simulate(s chess.ChessVariation, h *Node, depth int, p *POMCP) float64 {
	if math.Pow(gamma, float64(depth)) < epsilon {
		return 0
	}
	if len(h.children) == 0 { // expand
		h.children = create_all_children(s, h)
		return rollout(s, h, depth, p)
	}
	a := get_most_promising_action_by_ucb(s, h)
	o := s.ExecuteMove(a)
	possible_moves := o.PossibleMoves()
	if len(possible_moves) > 0 {
		o = o.ExecuteMove(random_element(possible_moves))
	}
	ha := h.children[Combined_Key{a, o.CreateView().String()}] // s needs to be an observation?
	if ha == nil {
		// if key does not exist -> create it
		h.children[Combined_Key{a, o.CreateView().String()}] = &Node{h, nil, 0, 0, a, o, nil}
		ha = h.children[Combined_Key{a, o.CreateView().String()}]
	}

	reward := simulate(o, ha, depth, p)
	h.beliefs = append(h.beliefs, s)
	h.visits++
	ha.visits++
	ha.value = ha.value + (reward-ha.value)/float64(ha.visits)
	return reward
}

func rollout(s chess.ChessVariation, h *Node, depth int, p *POMCP) float64 {
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
		return float64(result) // 1 for win -1 for lose
	}
	return rollout(new_s, h, depth+1, p)
}

func random_element[T any](collection []T) T {
	selected_index := 0
	size_collection := len(collection)
	if size_collection == 0 {
		var zero T
		return zero
	} else if size_collection > 1 {
		selected_index = rand.Intn(size_collection - 1)
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

func create_all_children(s chess.ChessVariation, h *Node) map[Combined_Key]*Node {
	possible_transitions := make(map[Combined_Key]*Node)
	// differentiate between own action and opponent action in data type
	for _, own_move := range s.PossibleMoves() {
		new_s := s.ExecuteMove(own_move)
		for _, opponent_move := range new_s.PossibleMoves() {
			new_s = new_s.ExecuteMove(opponent_move)
		}
		new_child := &Node{h, nil, 0, 0, own_move, new_s, nil}
		possible_transitions[Combined_Key{own_move, new_s.CreateView().String()}] = new_child
	}

	return possible_transitions
}

// simulates random own move and opponents
func state_transition(s chess.ChessVariation, h *Node) chess.ChessVariation {
	possible_moves := s.PossibleMoves()
	if len(possible_moves) == 0 {
		return s
	}
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
