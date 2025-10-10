package server

import (
	c "monte_carlo_hybrids/chess_variation"
	p "monte_carlo_hybrids/player"
)

func playGame(game c.ChessVariation, player1 p.Player, player2 p.Player) p.Player {
	move := 0
	var currPlayer p.Player
	game.InitGame()

	for !game.GameOver() {
		currPlayer = setCurrPlayer(move, player1, player2)
		currPlayer.GetMove()

		move++
	}

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
