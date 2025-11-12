package main

import (
	"monte_carlo_hybrids/chess_variation"
	"monte_carlo_hybrids/player"
	"monte_carlo_hybrids/server"
)

func main() {
	dark_pawn_chess := chess_variation.DarkPawnChess{}
	player2 := player.RandomPlayer{}
	player1 := player.POMCP{}
	greedy_wins := 0
	pomcp_wins := 0
	for i := 0; i < 10; i++ {
		winner := server.PlayGame(&dark_pawn_chess, &player1, &player2)
		if winner == &player1 {
			pomcp_wins++
		} else if winner == &player2 {
			greedy_wins++
		}
	}
	println("Greedy wins: ", greedy_wins)
	println("POMCP wins: ", pomcp_wins)
}
