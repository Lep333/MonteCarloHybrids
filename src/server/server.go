package server

import (
	"fmt"
	"monte_carlo_hybrids/chess_variation"
	c "monte_carlo_hybrids/chess_variation"
	"monte_carlo_hybrids/player"
	p "monte_carlo_hybrids/player"
)

func PlayGame(game c.ChessVariation, player1 p.Player, player2 p.Player) (int, []chess_variation.Move, []int) {
	move := 0
	var currPlayer p.Player
	game.InitGame()
	var move_history []chess_variation.Move
	var rollouts []int
	var result int
	for {
		game_over, tmp_result := game.GameOver()
		if game_over {
			result = tmp_result
			break
		}
		//fmt.Printf("Move: %v \n", move)
		//fmt.Printf("%v", game)
		white_to_play := move%2 == 0
		currPlayer = setCurrPlayer(move, player1, player2)
		moves := game.PossibleMoves()
		//fmt.Printf("Possible moves: %v \n", moves)
		if len(moves) == 0 {
			fmt.Printf("No possible moves anymore! \n")
			break
		}
		board := currPlayer.ViewFunc(game)
		move_to_play := currPlayer.GetMove(board, white_to_play)
		move_history = append(move_history, move_to_play)
		if pomcp, ok := currPlayer.(*player.POMCP); ok {
			rollouts = append(rollouts, pomcp.Rollouts)
		} else {
			rollouts = append(rollouts, 0)
		}
		fmt.Printf("Selected move: %v \n", move_to_play)
		// check if move legal!
		legal_move := false
		for _, move := range moves {
			if move == move_to_play {
				legal_move = true
				break
			}
		}
		if !legal_move {
			panic("---ILLEGAL MOVE---")
		}
		game = game.ExecuteMove(move_to_play)
		move++
	}
	fmt.Printf("%v", game)
	fmt.Printf("Game finished! \n")
	return result, move_history, rollouts
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
