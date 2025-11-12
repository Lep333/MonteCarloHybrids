package server

import (
	"fmt"
	c "monte_carlo_hybrids/chess_variation"
	p "monte_carlo_hybrids/player"
)

func PlayGame(game c.ChessVariation, player1 p.Player, player2 p.Player) p.Player {
	move := 0
	var currPlayer p.Player
	game.InitGame()

	for {
		game_over, _ := game.GameOver()
		if game_over {
			break
		}
		fmt.Printf("Move: %v \n", move)
		fmt.Printf("%v", game)
		white_to_play := move%2 == 0
		currPlayer = setCurrPlayer(move, player1, player2)
		moves := game.PossibleMoves()
		fmt.Printf("Possible moves: %v \n", moves)
		if len(moves) == 0 {
			fmt.Printf("No possible moves anymore! \n")
			break
		}
		board := game.CreateView()
		move_to_play := currPlayer.GetMove(board, white_to_play)
		fmt.Printf("Selected move: %v \n", move_to_play)
		// TODO: check if move legal!
		legal_move := false
		for _, move := range moves {
			if move == move_to_play {
				legal_move = true
				break
			}
		}
		if !legal_move {
			println("---ILLEGAL MOVE---")
			return nil
		}
		game = game.ExecuteMove(move_to_play)
		move++
	}
	fmt.Printf("%v", game)
	fmt.Printf("Game finished! \n")
	return currPlayer
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
