package main

import (
	"monte_carlo_hybrids/chess_variation"
	"monte_carlo_hybrids/player"
	"monte_carlo_hybrids/server"
)

func main() {
	dark_pawn_chess := chess_variation.DarkPawnChess{}
	player1 := player.RandomPlayer{}
	player2 := player.MonteCarlo{}
	server.PlayGame(&dark_pawn_chess, &player1, &player2)
}
