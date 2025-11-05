package player

import (
	"fmt"
	"math"
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
	"slices"
	"time"
)

type MonteCarlo struct {
	particle_filter []chess.ChessVariation
}

type MonteCarloTreeNode struct {
	parent   *MonteCarloTreeNode
	children []*MonteCarloTreeNode
	visits   uint
	wins     uint
	move     chess.Move
	board    chess.ChessVariation
}

func (m *MonteCarlo) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	// TODO: settings
	time_length_in_ms := 1000
	// particle_filter_size := 500

	start := time.Now()

	fill_particle_filter(m, board)
	root := MonteCarloTreeNode{nil, nil, 0, 0, chess.Move{}, board}
	pomcp(m, board, &root, start, time_length_in_ms)

	slices.SortFunc(root.children, func(a, b *MonteCarloTreeNode) int {
		return int(b.wins - a.wins)
	})
	fmt.Printf("Root visits: %v \n", root.visits)
	fmt.Printf("root %v \n", root.children)
	for i, child := range root.children {
		fmt.Printf("Child %v: wins %v visits: %v \n", i, child.wins, child.visits)
	}
	fmt.Printf("Particle Filter size: %v \n", len(m.particle_filter))
	return root.children[0].move
}

func fill_particle_filter(m *MonteCarlo, board chess.ChessVariation) {
	if m.particle_filter == nil {
		m.particle_filter = append(m.particle_filter, board)
	}
	if board.GetNumberOfMoves() == 0 {
		return
	}
	var new_particle_filter []chess.ChessVariation
	for _, state := range m.particle_filter {
		// TODO: pick state and move randomly
		// also limit filter size
		curr_view := state.CreateView()
		poss_moves := state.GetPreviousBoard().PossibleMoves()
		for _, poss_move := range poss_moves {
			new_board := state.GetPreviousBoard().ExecuteMove(poss_move)
			view := new_board.CreateView()
			// check if result fits with knowledge
			if view.String() == curr_view.String() {
				new_particle_filter = append(new_particle_filter, new_board)
			}
		}
	}
	upper_bound := min(len(new_particle_filter), 5000)
	m.particle_filter = new_particle_filter[0:upper_bound]
}

func pomcp(
	m *MonteCarlo,
	board chess.ChessVariation,
	root *MonteCarloTreeNode,
	start time.Time,
	time_length_in_ms int) {
	moves := board.PossibleMoves()
	for _, move := range moves {
		new_board := board.ExecuteMove(move)
		new_node := MonteCarloTreeNode{root, nil, 0, 0, move, new_board}
		root.children = append(root.children, &new_node)
	}

	for int(time.Since(start).Milliseconds()) < time_length_in_ms {
		state := m.particle_filter[rand.Intn(len(m.particle_filter)-1)]
		for _, child := range root.children {
			child.board = state.ExecuteMove(child.move)
		}

		curr_node := selection(root)
		expansion(curr_node)
		curr_board := rollout(curr_node)
		backpropagation(curr_node, curr_board, curr_board.GetNumberOfMoves()%2 == 0)
	}
}

func selection(root *MonteCarloTreeNode) *MonteCarloTreeNode {
	curr_node := root
	for len(curr_node.children) > 0 {
		ucb_sort(curr_node)
		curr_node = curr_node.children[0]
	}
	return curr_node
}

func expansion(curr_node *MonteCarloTreeNode) {
	poss_moves := curr_node.board.PossibleMoves()
	for _, move := range poss_moves {
		new_board := curr_node.board.ExecuteMove(move)
		new_board.CreateView()
		new_node := MonteCarloTreeNode{curr_node, nil, 0, 0, move, new_board}
		curr_node.children = append(curr_node.children, &new_node)
	}
}

func rollout(curr_node *MonteCarloTreeNode) chess.ChessVariation {
	var poss_next_state []chess.ChessVariation
	curr_board := curr_node.board
	for {
		game_over, _ := curr_board.GameOver()
		if game_over {
			break
		}
		poss_next_state = []chess.ChessVariation{}
		poss_moves := curr_board.PossibleMoves()

		for _, move := range poss_moves {
			new_board := curr_board.ExecuteMove(move)
			poss_next_state = append(poss_next_state, new_board)
		}
		rand.Shuffle(len(poss_next_state), func(i, j int) {
			poss_next_state[i], poss_next_state[j] = poss_next_state[j], poss_next_state[i]
		})
		if len(poss_next_state) > 0 {
			curr_board = poss_next_state[0]
		}
	}
	return curr_board
}

func backpropagation(curr_node *MonteCarloTreeNode, curr_board chess.ChessVariation, whiteToPlay bool) {
	_, result := curr_board.GameOver()
	if !whiteToPlay {
		if result == -1 {
			result = 1
		} else {
			result = 0
		}
	} else {
		if result != 1 {
			result = 0
		}
	}
	for curr_node != nil {
		curr_node.visits += 1
		curr_node.wins += uint(result)
		curr_node = curr_node.parent
	}
}

func ucb_sort(curr_node *MonteCarloTreeNode) {
	slices.SortFunc(curr_node.children, func(a, b *MonteCarloTreeNode) int {
		a_ucb := ucb(a)
		b_ucb := ucb(b)
		result := a_ucb - b_ucb
		if result > 0 {
			return -1
		} else {
			return 1
		}
	})
}

func ucb(a *MonteCarloTreeNode) float64 {
	if a.visits == 0 {
		return float64(math.MaxFloat64)
	}
	return float64(a.wins)/float64(a.visits) + math.Sqrt(1)*math.Sqrt(math.Log(float64(a.parent.visits)/float64(a.visits)))
}

func search(h *MonteCarloTreeNode) chess.Move {
	start := time.Now()
	duration_in_ms := 1000
	for (time.Since(start)) < time.Duration(duration_in_ms) {
		var s *MonteCarloTreeNode
		if h == nil {
			s = &MonteCarloTreeNode{}
		} else {
			s = h
		}
		simulate(s, h, 0)
	}
}

func simulate(s *MonteCarloTreeNode, h *MonteCarloTreeNode, depth uint) float64 {
	gamma := 0.95
	epsilon := 0.005
	if math.Pow(gamma, float64(depth)) < epsilon {
		return 0
	}
	if h.children == nil {
		var children []*MonteCarloTreeNode
		for _, move := range h.board.PossibleMoves() {
			children = append(children, &MonteCarloTreeNode{h.parent, nil, 0, 0, move, h.board})
		}
		h.children = children
		pomcp_rollout(s, h, depth+1)
	}
	ucb_sort(h)
	new_s := h.board.ExecuteMove(h.children[0].move)
	r := r + pomcp_rollout(s, h, depth+1)
}

func pomcp_rollout(s *MonteCarloTreeNode, h *MonteCarloTreeNode, depth uint) float64 {
	gamma := 0.95
	epsilon := 0.005
	if math.Pow(gamma, float64(depth)) < epsilon {
		return 0
	}
}
