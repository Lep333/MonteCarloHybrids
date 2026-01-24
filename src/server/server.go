package server

import (
	"fmt"
	"monte_carlo_hybrids/chess_variation"
	c "monte_carlo_hybrids/chess_variation"
	"monte_carlo_hybrids/player"
	p "monte_carlo_hybrids/player"
)

var poss_moves = [100]c.Move{}

func PlayGame(game c.ChessVariation, player1 p.Player, player2 p.Player) (
	int, []chess_variation.Move, []int, []int) {
	move := 0
	var currPlayer p.Player
	game.InitGame()
	var move_history []chess_variation.Move
	var rollouts []int
	var beliefs []int
	var result int
	for {
		game_over, tmp_result := game.GameOver()
		if game_over {
			result = tmp_result
			break
		}
		fmt.Printf("Move: %v \n", move)
		fmt.Printf("%v", game)
		white_to_play := move%2 == 0
		currPlayer = setCurrPlayer(move, player1, player2)
		n := game.PossibleMoves(poss_moves[:])
		fmt.Printf("Possible moves: %v \n", poss_moves[:n])
		if n == 0 {
			fmt.Printf("No possible moves anymore! \n")
			break
		}
		board := currPlayer.ViewFunc(game)
		move_to_play := currPlayer.GetMove(board, white_to_play)
		move_history = append(move_history, move_to_play)
		if pomcp, ok := currPlayer.(*player.POMCP); ok {
			rollouts = append(rollouts, pomcp.Rollouts)
			beliefs = append(beliefs, pomcp.NBeliefs)
		} else {
			rollouts = append(rollouts, 0)
			beliefs = append(beliefs, 0)
		}
		fmt.Printf("Selected move: %v \n", move_to_play)
		// check if move legal!
		legal_move := false
		println("n: ", n)
		for _, move := range poss_moves[:n] {
			if move == move_to_play {
				legal_move = true
				break
			}
		}
		if !legal_move {
			panic("---ILLEGAL MOVE---")
		}
		game.ExecuteMove(move_to_play)
		move++
	}
	fmt.Printf("%v", game)
	fmt.Printf("Game finished! \n")
	return result, move_history, rollouts, beliefs
}

func setCurrPlayer(move int, player1 p.Player, player2 p.Player) p.Player {
	var currPlayer p.Player

	if move%2 == 0 {
		currPlayer = player1
	} else {
		currPlayer = player2
	}
	return currPlayer
}
