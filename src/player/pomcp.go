package player

import (
	"math"
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
	"time"
)

const NODE_MAX = 100000
const BELIEFS_MAX = 1000

type POMCP struct {
	Root            *Node
	Started_playing bool
	Last_move       chess.Move
	Settings        Settings
	Rollouts        int
	NBeliefs        int
	Nodes           [NODE_MAX]Node
	node_count      uint
}

type Combined_Key struct {
	action      chess.Move
	observation string
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
	children    [1000]*Node
	child_count int16
	value       float64    // v
	visits      int        // n
	move        chess.Move // a
	observation uint64     // o
	beliefs     Belief     // B
}

type Belief struct {
	beliefs      [BELIEFS_MAX]chess.ChessVariation
	belief_count int64
}

//var Moves = [100]chess.Move{}

func (p *POMCP) get_node_address() *Node {
	// TODO: add kill flag
	node := &p.Nodes[p.node_count]
	p.node_count++
	if p.node_count == NODE_MAX {
		p.node_count = 0
	}
	return node
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
	p.NBeliefs = int(p.Root.beliefs.belief_count)
	if (selected_move == chess.Move{}) {
		selected_move = random_element(board.PossibleMoves())
	}
	p.Last_move = selected_move
	return selected_move
}

func (p *POMCP) Init_pomcp(board chess.ChessVariation, whiteToPlay bool) {
	if p.Root == nil {
		particle_filter := Belief{}
		if board.GetNumberOfMoves() == 0 {
			init_board := board.ReturnBoard()
			init_board.InitGame()
			particle_filter.beliefs[particle_filter.belief_count] = init_board
			particle_filter.belief_count++
			root := p.get_node_address()
			root.beliefs = particle_filter
			p.Root = root
		} else if board.GetNumberOfMoves() == 1 {
			init_board := board.ReturnBoard()
			init_board.InitGame()
			for _, move := range init_board.PossibleMoves() {
				game := init_board.ReturnBoard()
				game.ExecuteMove(move)
				particle_filter.beliefs[particle_filter.belief_count] = game
				particle_filter.belief_count++
			}
			root := p.get_node_address()
			root.beliefs = particle_filter
			p.Root = root
		}
		p.Started_playing = whiteToPlay
	}
}

func prune_tree_and_update_beliefs(p *POMCP, board chess.ChessVariation) {
	if p.Last_move != (chess.Move{}) {
		white := board.GetNumberOfMoves()%2 == 0
		board_hash := board.ViewHash(white)
		var new_root *Node
		for _, child := range p.Root.children {
			if child.move == p.Last_move && child.observation == board_hash {
				new_root = child
				break
			}
		}
		if new_root != nil {
			new_root.parent = nil
		}
		if new_root == nil {
			consistent_beliefs := Belief{}
			for _, belief := range p.Root.beliefs.beliefs[:p.Root.beliefs.belief_count] {
				new_s := p.state_transition(belief)
				if new_s.ViewHash(white) == board.ViewHash(white) {
					consistent_beliefs.beliefs[consistent_beliefs.belief_count] = belief
					consistent_beliefs.belief_count++
				}
			}
			root := p.get_node_address()
			root.beliefs = consistent_beliefs
			new_root = root
		}
		p.Root = new_root // TODO: find
		println(p.Root.beliefs.belief_count)
		// update consistent beliefs
		consistent_beliefs := Belief{}
		// consistent_beliefs[] = append(consistent_beliefs, board)
		for _, belief := range p.Root.beliefs.beliefs[:p.Root.beliefs.belief_count] {
			println(belief.ViewHash(white), board.ViewHash(white))
			if belief.ViewHash(white) == board.ViewHash(white) {
				consistent_beliefs.beliefs[consistent_beliefs.belief_count] = belief
				consistent_beliefs.belief_count++
			}
		}
		if consistent_beliefs.belief_count == 0 {
			consistent_beliefs.beliefs[0] = board
		}
		p.Root.beliefs = consistent_beliefs
	}
}

func (p *POMCP) search(h *Node) chess.Move {
	start_time := time.Now()
	//for i := 0; i < 50000; i++ {
	for !time_out(start_time, p.Settings.Termination_parameter) {
		s := random_belief(h.beliefs)
		copy := s.ReturnBoard()
		p.simulate(copy, h, 0, p.Settings.Gamma)
	}
	return get_best_move(h)
}

func (p *POMCP) simulate(s chess.ChessVariation, h *Node, depth int, discount float64) float64 {
	discount *= discount
	if discount < p.Settings.Epsilon {
		return 0
	}

	white := s.GetNumberOfMoves()%2 == 0
	if h.parent != nil {
		child := p.get_children(h.parent, h.move, h.observation)
		if child == nil { // expand
			h.parent.children[h.parent.child_count] = h
			h.parent.child_count++
			if p.Settings.Prior_hybrid != nil {
				return p.Settings.Prior_hybrid.Prior(h, s)
			}
			return p.rollout(s, depth, discount)
		}
	}
	var a chess.Move
	if over, _ := s.GameOver(); over {
		return p.rollout(s, depth, discount)
	}
	if p.Settings.Selection_hybrid != nil {
		a = p.Settings.Selection_hybrid.Select(s)
	} else {
		a = p.get_most_promising_action_by_ucb(s, h)
	}
	s.ExecuteMove(a)
	moves := s.PossibleMoves()
	if len(moves) > 0 {
		corrective := CorrectiveSelection{0.6, 0.05}
		opponent_move := corrective.Select(s)
		//opponent_move := random_element(s.PossibleMoves())
		s.ExecuteMove(opponent_move)
	}
	o := s.ViewHash(white)
	ha := p.get_children(h, a, o) // s needs to be an observation?
	if ha == nil {
		ha = p.get_node_address()
		ha.move = a
		ha.observation = o
		ha.parent = h
	}
	if depth == 0 {
		if (ha.beliefs == Belief{}) {
			ha.beliefs = Belief{}
		}
		got_belief_already := false
		for _, belief := range ha.beliefs.beliefs[:ha.beliefs.belief_count] {
			if belief == s {
				got_belief_already = true
				break
			}
		}
		if !got_belief_already {
			ha.beliefs.beliefs[ha.beliefs.belief_count] = s
			ha.beliefs.belief_count++
			if ha.beliefs.belief_count == BELIEFS_MAX {
				ha.beliefs.belief_count = 0
			}
		}
	}
	reward := p.Settings.Gamma * p.simulate(s, ha, depth+1, discount)
	h.visits++
	ha.visits++
	ha.value = ha.value + (reward-ha.value)/float64(ha.visits)
	return reward
}

func (p *POMCP) get_children(n *Node, a chess.Move, o uint64) *Node {
	var child *Node

	for _, node := range n.children[:n.child_count] {
		if node.move == a && node.observation == o {
			return node
		}
	}
	return child
}

func (p *POMCP) rollout(s chess.ChessVariation, depth int, discount float64) float64 {
	discount *= discount
	if discount < p.Settings.Epsilon {
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
	return p.Settings.Gamma * p.rollout(new_s, depth+1, discount)
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

func random_belief(b Belief) chess.ChessVariation {
	i := 0
	if b.belief_count > 1 {
		i = rand.Intn(int(b.belief_count))
	}
	for _, value := range b.beliefs[:b.belief_count] {
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
	for _, child := range h.children[:h.child_count] {
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

// simulates random own move and opponents
func (p *POMCP) state_transition(s chess.ChessVariation) chess.ChessVariation {
	moves := s.PossibleMoves()
	if len(moves) == 0 {
		return s
	}
	selected_move := p.rollout_move_selection(s, moves)
	s.ExecuteMove(selected_move)
	game_over, _ := s.GameOver()
	if !game_over {
		moves = s.PossibleMoves()
		opponent_selected_move := p.rollout_move_selection(s, moves)
		s.ExecuteMove(opponent_selected_move)
	}
	return s
}

func (p *POMCP) rollout_move_selection(s chess.ChessVariation, possible_moves []chess.Move) chess.Move {
	var capture_moves [1000]chess.Move
	var capture_moves_count uint

	for _, move := range possible_moves {
		if move.Capture {
			capture_moves[capture_moves_count] = move
			capture_moves_count++
		}
	}
	var selected_move chess.Move
	if len(capture_moves) > 0 && rand.Float64() < p.Settings.Rollout_capture {
		selected_move = random_element(capture_moves[:capture_moves_count])
	} else {
		selected_move = random_element(possible_moves)
	}
	return selected_move
}

var not_visited_action = [1000]chess.Move{}
var not_visited_count = 0

func (p *POMCP) get_most_promising_action_by_ucb(s chess.ChessVariation, h *Node) chess.Move {
	var max_ucb float64 = math.Inf(-1)
	var max_child *Node
	not_visited_count = 0
	not_visited := false
	white := s.GetNumberOfMoves()%2 == 0
	for _, a := range s.PossibleMoves() {
		s.ExecuteMove(a)
		o := s.ViewHash(white)
		s.UndoMove(a)
		var ucb float64
		child := p.get_children(h, a, o)
		if child != nil {
			if child.visits == 0 {
				ucb = math.MaxFloat64
				not_visited = true
				not_visited_action[not_visited_count] = child.move
				not_visited_count++
			} else {
				exploration := p.Settings.Ucb_c * math.Sqrt(math.Log(float64(h.visits))/float64(child.visits))
				ucb = child.value + exploration
			}
			if ucb > max_ucb {
				max_ucb = ucb
				max_child = child
			}
		} else {
			not_visited = true
			not_visited_action[not_visited_count] = a
			not_visited_count++
		}
	}

	if not_visited {
		return random_element(not_visited_action[:not_visited_count])
	}
	return max_child.move
}

func (p *POMCP) String() string {
	if p.Settings.POMCP_name == "" {
		return "POMCP"
	} else {
		return p.Settings.POMCP_name
	}
}
