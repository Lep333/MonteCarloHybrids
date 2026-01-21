package player

import (
	"fmt"
	"math"
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
	"strconv"
	"time"
)

type POMCP struct {
	Root            *Node
	Started_playing bool
	Last_move       chess.Move
	Settings        Settings
	Rollouts        int
	NBeliefs        int
	Nodes           map[string]*Node
}

type Combined_Key struct {
	action      chess.Move
	observation string
}

type History struct {
	history []Combined_Key
}

type Settings struct {
	Termination_parameter     int
	Gamma                     float64
	Epsilon                   float64
	Ucb_c                     float64
	Rollout_capture           float64
	Termination_func_param    float64
	POMCP_name                string
	Prior_hybrid              PriorHybrid
	Selection_hybrid          SelectionHybrid
	Rollout_selection         SelectionHybrid
	Early_playout_termination EarlyTermination
}

type Node struct {
	parent      *Node
	children    map[Combined_Key]*Node
	value       float64                         // v
	visits      int                             // n
	move        chess.Move                      // a
	observation chess.ChessVariation            // o
	beliefs     map[string]chess.ChessVariation // B
	history     string
}

func (p *POMCP) ViewFunc(board chess.ChessVariation) chess.ChessVariation {
	white := board.GetNumberOfMoves()%2 == 0
	return board.CreateView(white)
}

func (p *POMCP) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	if board.GetNumberOfMoves() < 2 {
		p.Root = nil
		p.Last_move = chess.Move{}
	}
	p.Init_pomcp(board, whiteToPlay)
	prune_tree_and_update_beliefs(p, board)
	visits_before := p.Root.visits
	selected_move := p.search(p.Root)
	p.Rollouts = p.Root.visits - visits_before
	p.NBeliefs = len(p.Root.beliefs)
	if (selected_move == chess.Move{}) {
		selected_move = random_element(board.PossibleMoves())
	}
	p.Last_move = selected_move
	return selected_move
}

func (p *POMCP) Init_pomcp(board chess.ChessVariation, whiteToPlay bool) {
	if p.Root == nil {
		white := board.GetNumberOfMoves()%2 == 0
		beliefs := make(map[string]chess.ChessVariation)
		if board.GetNumberOfMoves() == 0 {
			init_board := board.ReturnBoard()
			init_board.InitGame()
			beliefs[init_board.String()] = init_board
			p.Root = &Node{nil, map[Combined_Key]*Node{}, 0, 0, chess.Move{}, board.CreateView(white), beliefs, ""}
		} else if board.GetNumberOfMoves() == 1 {
			init_board := board.ReturnBoard()
			init_board.InitGame()
			for _, move := range init_board.PossibleMoves() {
				game := init_board.ReturnBoard()
				game.ExecuteMove(move)
				follow_state_string := game.String()
				if _, belief := beliefs[follow_state_string]; !belief {
					beliefs[follow_state_string] = game
				}
			}
			p.Root = &Node{nil, map[Combined_Key]*Node{}, 0, 0, chess.Move{}, nil, beliefs, ""}
		}
		p.Started_playing = whiteToPlay
		p.Nodes = map[string]*Node{}
		p.Nodes[""] = p.Root
	}
}

func prune_tree_and_update_beliefs(p *POMCP, board chess.ChessVariation) {
	if p.Last_move != (chess.Move{}) {
		white := board.GetNumberOfMoves()%2 == 0
		observation_string := board.CreateView(white).String()
		new_history := p.Root.history + fmt.Sprintf("<%s,%s>", moveToStr(p.Last_move), observation_string)
		new_root := p.Nodes[new_history]
		if new_root != nil {
			new_root.parent = nil
		}
		if new_root == nil {
			consistent_beliefs := map[string]chess.ChessVariation{}
			for key, belief := range p.Root.beliefs {
				new_s := p.state_transition(belief)

				if new_s.CreateView(white).String() == board.String() {
					consistent_beliefs[key] = belief
				}
			}
			new_root = &Node{nil, map[Combined_Key]*Node{}, 0, 0, chess.Move{}, nil, consistent_beliefs, new_history}
		}
		p.Root = new_root // TODO: find
		// update consistent beliefs
		consistent_beliefs := map[string]chess.ChessVariation{}
		// consistent_beliefs[] = append(consistent_beliefs, board)
		for key, belief := range p.Root.beliefs {
			if belief.CreateView(white).String() == board.String() {
				consistent_beliefs[key] = belief
			}
		}
		if len(consistent_beliefs) == 0 {
			consistent_beliefs[board.String()] = board
		}
		p.Root.beliefs = consistent_beliefs
	}
}

func (p *POMCP) search(h *Node) chess.Move {
	start_time := time.Now()
	// for i := 0; i < 10000; i++ {
	for !time_out(start_time, p.Settings.Termination_parameter) {
		s := random_belief(h.beliefs)
		copy := s.ReturnBoard()
		p.simulate(copy, h, 0, p.Settings.Gamma)
	}
	return get_best_move(h)
}

func (p *POMCP) simulate(s chess.ChessVariation, h *Node, depth int, gamma float64) float64 {
	if gamma < p.Settings.Epsilon {
		return 0
	}
	own_action := depth%2 == 0
	var opp bool
	opp = false
	if h.parent != nil {
		_, opp = h.parent.children[Combined_Key{h.move, h.observation.String()}]
		opp = !opp
	}
	expansion := opp
	if own_action && depth > 1 {
		_, own := p.Nodes[h.history]
		expansion = !own
	}
	white := s.GetNumberOfMoves()%2 == 0
	if expansion {
		// h.beliefs = append(h.beliefs, s)
		// h.children = create_all_children(s, h)
		if own_action && depth > 1 {
			p.Nodes[h.history] = h
		} else {
			h.parent.children[Combined_Key{h.move, h.observation.String()}] = h
		}
		if p.Settings.Prior_hybrid != nil {
			return p.Settings.Prior_hybrid.Prior(h, s)
		}
		return p.rollout(s, depth)
	}
	var a chess.Move
	if p.Settings.Selection_hybrid != nil {
		a = p.Settings.Selection_hybrid.Select(s)
	} else {
		a = p.get_most_promising_action_by_ucb(s, h)
	}
	h.beliefs[s.String()] = s.ReturnBoard()
	s.ExecuteMove(a)
	o := s.CreateView(white)
	var ha *Node
	if own_action {
		ha = h.children[Combined_Key{a, o.String()}]
		if ha == nil {
			// if key does not exist -> create it
			ha = &Node{
				h, map[Combined_Key]*Node{}, 0, 0, a, o, map[string]chess.ChessVariation{}, h.history,
			}
		}
	} else {
		observation_string := h.observation.String()
		new_history := h.history + fmt.Sprintf("<%s,%s>", moveToStr(h.move), observation_string)
		ha = p.Nodes[new_history]
		if ha == nil {
			// if key does not exist -> create it
			ha = &Node{
				h, map[Combined_Key]*Node{}, 0, 0, h.move, o, map[string]chess.ChessVariation{}, new_history,
			}
		}
		ha.beliefs[s.String()] = s.ReturnBoard()
	}

	reward := p.Settings.Gamma * p.simulate(s, ha, depth+1, math.Pow(gamma, float64(depth)))
	h.visits++
	ha.visits++
	ha.value = ha.value + (reward-ha.value)/float64(ha.visits)
	return -reward
}

func (p *POMCP) rollout(s chess.ChessVariation, depth int) float64 {
	if math.Pow(p.Settings.Gamma, float64(depth)) < p.Settings.Epsilon {
		return 0
	}

	// early playout termination?
	if p.Settings.Early_playout_termination != nil {
		termination, value := p.Settings.Early_playout_termination.
			EarlyPlayoutTermination(s, float64(depth))
		if termination {
			return value
		}
	}

	var new_s chess.ChessVariation
	if p.Settings.Rollout_selection != nil {
		move := p.Settings.Rollout_selection.Select(s)
		s.ExecuteMove(move)
		new_s = s
	} else {
		new_s = p.state_transition(s)
	}

	game_over, result := new_s.GameOver()
	if game_over {
		if !p.Started_playing {
			if result == 1 {
				result = -1
			} else {
				result = 1
			}
		}
		return float64(result) // 1 for win -1 for lose
	}
	return p.Settings.Gamma * p.rollout(new_s, depth+1)
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
	i := 0
	if len(m) > 1 {
		i = rand.Intn(len(m))
	}
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
			visits := float64(child.visits)
			if child.visits == 0 {
				visits = 1
			}
			accumulated_action[child.move] = [2]float64{(val[0]*val[1] + float64(child.visits)*child.value) / (val[1] + visits), val[1] + float64(child.visits)}
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
	white := s.GetNumberOfMoves()%2 == 0
	// differentiate between own action and opponent action in data type
	for _, belief := range h.beliefs {
		for _, own_move := range belief.PossibleMoves() {
			new_s := belief.ReturnBoard()
			new_s.ExecuteMove(own_move)
			belief := map[string]chess.ChessVariation{}
			belief[new_s.String()] = new_s
			o := new_s.CreateView(white)
			new_child := &Node{h, map[Combined_Key]*Node{}, 0, 0, own_move, o, belief, ""}
			possible_transitions[Combined_Key{own_move, o.String()}] = new_child
		}
	}
	return possible_transitions
}

// simulates random own move
func (p *POMCP) state_transition(s chess.ChessVariation) chess.ChessVariation {
	possible_moves := s.PossibleMoves()
	if len(possible_moves) == 0 {
		return s
	}
	selected_move := p.rollout_move_selection(s, possible_moves)
	s.ExecuteMove(selected_move)
	return s
}

func (p *POMCP) rollout_move_selection(s chess.ChessVariation, possible_moves []chess.Move) chess.Move {
	var capture_moves []chess.Move
	for _, move := range possible_moves {
		if move.Capture {
			capture_moves = append(capture_moves, move)
		}
	}
	var selected_move chess.Move
	if len(capture_moves) > 0 && rand.Float64() < p.Settings.Rollout_capture {
		selected_move = random_element(capture_moves)
	} else {
		selected_move = random_element(possible_moves)
	}
	return selected_move
}

func (p *POMCP) get_most_promising_action_by_ucb(s chess.ChessVariation, h *Node) chess.Move {
	var max_ucb float64 = math.Inf(-1)
	var max_child *Node
	var not_visited_action []chess.Move
	not_visited := false

	if len(h.children) == 0 {
		return random_element(s.PossibleMoves())
	}
	for _, child := range h.children {
		var ucb float64
		if child.visits == 0 {
			ucb = math.MaxFloat64
			not_visited = true
			not_visited_action = append(not_visited_action, child.move)
		} else {
			exploration := p.Settings.Ucb_c * math.Sqrt(math.Log(float64(h.visits))/float64(child.visits))
			ucb = child.value + exploration
		}
		if ucb > max_ucb {
			max_ucb = ucb
			max_child = child
		}
	}
	if not_visited {
		return random_element(not_visited_action)
	}
	return max_child.move
}

func moveToStr(move chess.Move) string {
	return strconv.Itoa(int(move.From)) + "-" + strconv.Itoa(int(move.To))
}

func (p *POMCP) String() string {
	if p.Settings.POMCP_name == "" {
		return "POMCP"
	} else {
		return p.Settings.POMCP_name
	}
}
