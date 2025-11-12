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
	value       float64                         // v
	visits      int                             // n
	move        chess.Move                      // a
	observation chess.ChessVariation            // o
	beliefs     map[string]chess.ChessVariation // B
}

func (p *POMCP) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	time_limit := 1000
	if board.GetNumberOfMoves() < 2 {
		p.root = nil
		p.last_move = chess.Move{}
	}
	init_pomcp(board, p, whiteToPlay)
	prune_tree_and_update_beliefs(p, board)
	selected_move := search(p.root, time_limit, p)
	println("visits: ", p.root.visits)
	p.last_move = selected_move
	return selected_move
}

func init_pomcp(board chess.ChessVariation, p *POMCP, whiteToPlay bool) {
	if p.root == nil {
		beliefs := make(map[string]chess.ChessVariation)
		if board.GetNumberOfMoves() == 0 {
			init_board := board.ReturnBoard()
			init_board.InitGame()
			beliefs[init_board.String()] = init_board
			p.root = &Node{nil, nil, 0, 0, chess.Move{}, nil, beliefs}
		} else if board.GetNumberOfMoves() == 1 {
			init_board := board.ReturnBoard()
			init_board.InitGame()
			for _, move := range init_board.PossibleMoves() { // TODO: board with state 0
				follow_state := init_board.ExecuteMove(move)
				follow_state_string := follow_state.String()
				if _, belief := beliefs[follow_state_string]; !belief {
					beliefs[follow_state_string] = follow_state
				}
			}
			p.root = &Node{nil, nil, 0, 0, chess.Move{}, nil, beliefs}
		}
		p.started_playing = whiteToPlay
	}
}

func prune_tree_and_update_beliefs(p *POMCP, board chess.ChessVariation) {
	if p.last_move != (chess.Move{}) {
		board_string := board.String()
		new_root := p.root.children[Combined_Key{p.last_move, board_string}]
		if new_root == nil {
			consistent_beliefs := map[string]chess.ChessVariation{}
			for key, belief := range p.root.beliefs {
				new_s := state_transition(belief, p.root)
				if new_s.CreateView().String() == board.String() {
					consistent_beliefs[key] = belief
				}
			}
			new_root = &Node{nil, nil, 0, 0, chess.Move{}, nil, consistent_beliefs}
		}
		p.root = new_root // TODO: find
		// update consistent beliefs
		consistent_beliefs := map[string]chess.ChessVariation{}
		// consistent_beliefs[] = append(consistent_beliefs, board)
		for key, belief := range p.root.beliefs {
			if belief.CreateView().String() == board.String() {
				consistent_beliefs[key] = belief
			}
		}
		p.root.beliefs = consistent_beliefs
	}
}

func search(h *Node, time_limit int, p *POMCP) chess.Move {
	start_time := time.Now()
	// for i := 0; i < 10000; i++ {
	for !time_out(start_time, time_limit) {
		s := random_belief(h.beliefs)
		simulate(s, h, 0, p)
	}
	return get_best_move(h)
}

func simulate(s chess.ChessVariation, h *Node, depth int, p *POMCP) float64 {
	if math.Pow(gamma, float64(depth)) < epsilon {
		return 0
	}
	if len(h.children) == 0 { // expand
		// h.beliefs = append(h.beliefs, s)
		h.children = create_all_children(s, h)
		return rollout(s, h, depth, p)
	}
	a := get_most_promising_action_by_ucb(h)
	o := s.ExecuteMove(a)
	possible_moves := o.PossibleMoves()
	if len(possible_moves) > 0 {
		o = o.ExecuteMove(random_element(possible_moves))
	}
	ha := h.children[Combined_Key{a, o.CreateView().String()}] // s needs to be an observation?
	if ha == nil {
		// if key does not exist -> create it
		h.children[Combined_Key{a, o.CreateView().String()}] = &Node{h, nil, 0, 0, a, o, map[string]chess.ChessVariation{}}
		ha = h.children[Combined_Key{a, o.CreateView().String()}]
	}
	reward := simulate(o, ha, depth+1, p)
	h.beliefs[s.String()] = s
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
		selected_index = rand.Intn(size_collection)
	}
	return collection[selected_index]
}

func random_belief(m map[string]chess.ChessVariation) chess.ChessVariation {
	i := rand.Intn(len(m))
	for _, value := range m {
		if i == 0 {
			return value
		}
		i--
	}
	panic("empty belief set")
}

func time_out(start time.Time, time_limit int) bool {
	time_out := false
	if time.Since(start).Milliseconds() > int64(time_limit) {
		time_out = true
	}
	return time_out
}

func get_best_move(h *Node) chess.Move {
	accumulated_action := make(map[chess.Move][2]float64)
	var best_action chess.Move
	best_action_value := math.Inf(-1)
	for _, child := range h.children {
		if val, ok := accumulated_action[child.move]; ok {
			accumulated_action[child.move] = [2]float64{(val[0]*val[1] + float64(child.visits)*child.value) / (val[1] + float64(child.visits)), val[1] + float64(child.visits)}
		} else {
			accumulated_action[child.move] = [2]float64{child.value, float64(child.visits)}
		}
	}
	for action, value := range accumulated_action {
		if value[0] > best_action_value {
			best_action = action
			best_action_value = value[0]
		}
	}
	return best_action
}

func create_all_children(s chess.ChessVariation, h *Node) map[Combined_Key]*Node {
	possible_transitions := make(map[Combined_Key]*Node)
	// differentiate between own action and opponent action in data type
	for _, own_move := range s.PossibleMoves() {
		new_s := s.ExecuteMove(own_move)
		game_over, _ := new_s.GameOver()
		if game_over {
			belief := map[string]chess.ChessVariation{}
			belief[new_s.String()] = new_s
			new_child := &Node{h, nil, 0, 0, own_move, new_s, belief}
			possible_transitions[Combined_Key{own_move, new_s.CreateView().String()}] = new_child
		}
		for _, opponent_move := range new_s.PossibleMoves() {
			o := new_s.ExecuteMove(opponent_move)
			// TODO: what if game already over?
			if val, ok := possible_transitions[Combined_Key{own_move, o.CreateView().String()}]; ok {
				val.beliefs[o.String()] = o
			} else {
				belief := map[string]chess.ChessVariation{}
				belief[o.String()] = o
				new_child := &Node{h, nil, 0, 0, own_move, o, belief}
				possible_transitions[Combined_Key{own_move, o.CreateView().String()}] = new_child
			}
		}
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

func get_most_promising_action_by_ucb(h *Node) chess.Move {
	var max_ucb float64 = 0
	var max_child *Node
	var c float64 = 1

	for _, child := range h.children {
		var ucb float64
		if child.visits == 0 {
			ucb = math.MaxFloat64
		} else {
			ucb = child.value + c*math.Sqrt(math.Log(float64(h.visits))/float64(child.visits))
		}
		if max_child == nil || ucb > max_ucb {
			max_ucb = ucb
			max_child = child
		}
	}

	return max_child.move
}
