package player

import (
	"fmt"
	"math/rand"
	chess "monte_carlo_hybrids/chess_variation"
	"slices"
	"time"
)

type MonteCarlo struct{}

type MonteCarloTreeNode struct {
	parent   *MonteCarloTreeNode
	children []*MonteCarloTreeNode
	visits   uint
	wins     uint
	move     chess.Move
	board    chess.ChessVariation
}

func (m *MonteCarlo) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	time_length_in_ms := 500
	start := time.Now()
	root := MonteCarloTreeNode{nil, nil, 0, 0, chess.Move{}, board}
	moves := board.PossibleMoves()
	for _, move := range moves {
		new_board := board.ExecuteMove(move)
		new_node := MonteCarloTreeNode{&root, nil, 0, 0, move, new_board}
		root.children = append(root.children, &new_node)
	}

	for int(time.Now().Sub(start).Milliseconds()) < time_length_in_ms {
		curr_node := selection(&root)
		expansion(curr_node)
		curr_board := rollout(curr_node)
		backpropagation(curr_node, curr_board, whiteToPlay)
	}
	slices.SortFunc(root.children, func(a, b *MonteCarloTreeNode) int {
		return int(b.wins - a.wins)
	})
	fmt.Printf("Root visits: %v \n", root.visits)
	return root.children[0].move
}

func selection(root *MonteCarloTreeNode) *MonteCarloTreeNode {
	curr_node := root
	for len(curr_node.children) > 0 {
		rand.Shuffle(len(curr_node.children), func(i, j int) {
			curr_node.children[i], curr_node.children[j] = curr_node.children[j], curr_node.children[i]
		})
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
