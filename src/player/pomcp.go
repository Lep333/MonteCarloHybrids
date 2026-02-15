package player

import (
	"math"
	"math/rand"
	"monte_carlo_hybrids/chess_variation"
	chess "monte_carlo_hybrids/chess_variation"
	"time"
)

const NODE_MAX = 200000
const BELIEFS_MAX = 10000

type POMCP struct {
	Root            *Node
	Started_playing bool
	Last_move       chess.Move
	Settings        Settings
	Rollouts        int
	NBeliefs        int
	Nodes           [NODE_MAX]Node
	node_count      uint
	beliefs         *Belief
	child_pool      [NODE_MAX][200]*Node
	stack           []uint
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
	Opponent_modelling        bool
	OM_Threshold              float64
}

type Node struct {
	parent      *Node
	children    *[200]*Node
	value       float64    // v
	visits      int        // n
	observation uint64     // o
	move        chess.Move // a
	child_count int16
	free        bool
	expanded    bool
}

type Belief struct {
	beliefs      map[uint64]chess.ChessVariation
	belief_count int64
	allocated    bool
}

func (p *POMCP) get_node_address() *Node {
	found := false
	counter := NODE_MAX
	var node *Node
	var temp_node *Node
	for {
		if found {
			break
		}
		temp_node = &p.Nodes[p.node_count]
		if temp_node.free {
			found = true
			temp_node.children = &p.child_pool[p.node_count]
			temp_node.child_count = 0
			temp_node.free = false
		}
		p.node_count++
		if p.node_count == NODE_MAX {
			p.node_count = 0
		}
		if counter == 0 {
			break
		}
		counter--
	}
	if p.node_count == NODE_MAX {
		p.node_count = 0
	}
	if found {
		node = temp_node
	}
	return node
}

func (p *POMCP) ViewFunc(board chess.ChessVariation) chess.ChessVariation {
	white := p.Started_playing
	return board.CreateView(white)
}

func (p *POMCP) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	if board.GetNumberOfMoves() < 2 {
		p.Root = nil
		p.Last_move = chess.Move{}
	}
	p.Init_pomcp(board, whiteToPlay)
	p.prune_tree_and_update_beliefs(board)

	selected_move := p.search(p.Root)
	p.NBeliefs = len(p.beliefs.beliefs)
	p.Last_move = selected_move

	return selected_move
}

func (p *POMCP) Init_pomcp(board chess.ChessVariation, whiteToPlay bool) {
	if p.Root == nil {
		p.Nodes = [NODE_MAX]Node{}
		p.beliefs = &Belief{}
		p.beliefs.beliefs = make(map[uint64]chess.ChessVariation, 0)
		p.node_count = 0
		p.free_nodes()
		if board.GetNumberOfMoves() == 0 {
			init_board := board.ReturnBoard()
			init_board.InitGame()
			root := p.get_node_address()
			p.beliefs.beliefs[init_board.Hash()] = init_board
			p.beliefs.belief_count++
			p.Root = root
		} else if board.GetNumberOfMoves() == 1 {
			init_board := board.ReturnBoard()
			init_board.InitGame()
			root := p.get_node_address()
			for _, move := range init_board.PossibleMoves() {
				game := init_board.ReturnBoard()
				game.ExecuteMove(move)
				p.beliefs.beliefs[game.Hash()] = game
				p.beliefs.belief_count++
			}
			p.Root = root
		}
		p.Root.expanded = true
		p.Started_playing = whiteToPlay
	}
}

func (p *POMCP) prune_tree_and_update_beliefs(board chess.ChessVariation) {
	if p.Last_move != (chess.Move{}) {
		white := p.Started_playing
		board_hash := board.ViewHash(white)
		new_root := p.get_children(p.Root, p.Last_move, board_hash)
		if new_root == nil {
			root := p.get_node_address()
			if root == nil {
				root = &p.Nodes[0]
				root.child_count = 0
			}
			new_root = root
		}
		// update consistent beliefs
		beliefs := &Belief{}
		beliefs.beliefs = make(map[uint64]chess.ChessVariation)
		// consistent_beliefs[] = append(consistent_beliefs, board)
		for i := 0; i < 10000; i++ {
			belief := random_belief(p.beliefs.beliefs)
			copy := belief.ReturnBoard()
			copy.ExecuteMove(p.Last_move)
			opponent_move := random_element(copy.PossibleMoves())
			copy.ExecuteMove(opponent_move)
			bel_hash := copy.ViewHash(white)
			if bel_hash == board_hash {
				copy_hash := copy.Hash()
				if _, ok := beliefs.beliefs[copy_hash]; !ok {
					beliefs.beliefs[copy_hash] = copy
				}
			}
		}
		// create fallback particle
		if len(beliefs.beliefs) == 0 {
			for i, belief := range p.beliefs.beliefs {
				consistent_particle := board.Create_fallback_particle(belief, p.Started_playing)
				beliefs.beliefs[consistent_particle.Hash()] = consistent_particle
				if i == 100 {
					break
				}
			}
		}
		p.Root = new_root
		p.free_nodes()
		p.beliefs = beliefs
	}
}

func (p *POMCP) free_nodes() {
	for i := 0; i < NODE_MAX; i++ {
		p.Nodes[i].free = true
	}
	if p.Root != nil {
		lock_nodes(p.Root)
	}
}

func lock_nodes(node *Node) {
	node.free = false
	for _, child := range node.children[:node.child_count] {
		lock_nodes(child)
	}
}

func (p *POMCP) search(h *Node) chess.Move {
	start_time := time.Now()
	counter := 0
	//for i := 0; i < 10000; i++ {
	for !time_out(start_time, p.Settings.Termination_parameter) {
		s := random_belief(p.beliefs.beliefs)
		copy := s.ReturnBoard()
		p.simulate(copy, h, 0, p.Settings.Gamma)
		counter++
	}
	print("rollouts: ", counter, "beliefs: ", len(p.beliefs.beliefs))
	p.Rollouts = counter
	return get_best_move(h)
}

func (p *POMCP) simulate(s chess.ChessVariation, h *Node, depth int, current_gamma float64) float64 {
	if depth > 100 {
		return 0
	}

	white := p.Started_playing
	if !h.expanded {
		h.expanded = true
		if p.Settings.Prior_hybrid != nil {
			return p.Settings.Prior_hybrid.Prior(h, s)
		}
		return p.rollout(s, depth, current_gamma)
	}
	var a chess.Move
	if over, _ := s.GameOver(); over {
		return p.rollout(s, depth, current_gamma)
	}
	if p.Settings.Selection_hybrid != nil {
		a = p.Settings.Selection_hybrid.Select(s)
	} else {
		a = p.get_most_promising_action_by_ucb(s, h)
	}
	s.ExecuteMove(a)
	moves := s.PossibleMoves()
	p.opponent_modelling(s, moves)
	o := s.ViewHash(white)
	ha := p.get_children(h, a, o)
	if ha == nil {
		ha = p.create_node(h, a, o)
		if ha == nil {
			// start rollout if no node available
			return p.rollout(s, depth, current_gamma)
		}
	}
	reward := p.Settings.Gamma * p.simulate(s, ha, depth+2, current_gamma*p.Settings.Gamma)

	h.visits++
	ha.visits++
	ha.value = ha.value + (reward-ha.value)/float64(ha.visits)
	return reward
}

func (p *POMCP) opponent_modelling(s chess_variation.ChessVariation, moves []chess.Move) {
	if len(moves) > 0 {
		var opponent_move chess.Move
		if p.Settings.Opponent_modelling {
			greedy := GreedySelection{p.Settings.OM_Threshold}
			opponent_move = greedy.Select(s)
		} else {
			opponent_move = random_element(s.PossibleMoves())
		}
		s.ExecuteMove(opponent_move)
	}
}

func (p *POMCP) create_node(h *Node, a chess.Move, o uint64) *Node {
	ha := p.get_node_address()
	if ha == nil {
		return nil
	}
	ha.move = a
	ha.observation = o
	ha.parent = h
	ha.visits = 0
	ha.value = 0
	ha.free = false
	if h.child_count < 200 {
		h.children[h.child_count] = ha
		h.child_count++
	}
	return ha
}

func (p *POMCP) get_children(n *Node, a chess.Move, o uint64) *Node {
	var child *Node

	if n.child_count == 0 {
		return child
	}
	for _, node := range n.children[:n.child_count] {
		if node.move == a && node.observation == o {
			return node
		}
	}
	return child
}

func (p *POMCP) rollout(s chess.ChessVariation, depth int, current_gamma float64) float64 {
	if depth > 100 {
		return 0
	}
	// early playout termination?
	if p.Settings.Early_playout_termination != nil {
		termination, value := p.Settings.Early_playout_termination.
			EarlyPlayoutTermination(s, float64(depth), p.Started_playing)
		if termination {
			return value
		}
	}

	var new_s chess.ChessVariation
	if p.Settings.Rollout_selection != nil {
		move := p.Settings.Rollout_selection.Select(s)
		s.ExecuteMove(move)
		new_s = s
		game_over, _ := s.GameOver()
		if !game_over {
			opponent_selected_move := p.Settings.Rollout_selection.Select(s)
			s.ExecuteMove(opponent_selected_move)
		}
	} else {
		new_s = p.state_transition(s)
	}

	game_over, result := new_s.GameOver()
	if game_over {
		if !p.Started_playing {
			result = -result
		}
		return float64(result) // 1 for win -1 for lose
	}
	return p.Settings.Gamma * p.rollout(new_s, depth+2, current_gamma*p.Settings.Gamma)
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

func random_belief(b map[uint64]chess.ChessVariation) chess.ChessVariation {
	i := 0
	no_beliefs := len(b)
	if no_beliefs > 1 {
		i = rand.Intn(no_beliefs)
	}
	for _, value := range b {
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
	if capture_moves_count > 0 && rand.Float64() < p.Settings.Rollout_capture {
		selected_move = random_element(capture_moves[:capture_moves_count])
	} else {
		selected_move = random_element(possible_moves)
	}
	return selected_move
}

func (p *POMCP) get_most_promising_action_by_ucb(s chess.ChessVariation, h *Node) chess.Move {
	var not_visited_action = [1000]chess.Move{}
	var not_visited_count = 0
	var max_ucb float64 = math.Inf(-1)
	var max_child *Node
	not_visited := false
	white := p.Started_playing
	for _, a := range s.PossibleMoves() {
		copy := s.ReturnBoard()
		copy.ExecuteMove(a)
		o := copy.ViewHash(white)
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
	if max_child == nil {
		return random_element(s.PossibleMoves())
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
