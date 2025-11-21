package main

import (
	"fmt"
	"monte_carlo_hybrids/chess_variation"
	"monte_carlo_hybrids/player"
	"monte_carlo_hybrids/server"
	"os"
)

func main() {
	dark_pawn_chess := chess_variation.DarkPawnChess{}
	var player1, player2 player.Player
	player1 = &player.RandomPlayer{}
	settings := player.Settings{
		Termination_parameter: 2000,
		Gamma:                 0.95,
		Epsilon:               0.005,
		Ucb_c:                 8,
		Capture_reward:        0.2,
		Rollout_capture:       0.95,
	}
	player2 = &player.POMCP{Root: nil, Started_playing: false, Last_move: chess_variation.Move{}, Settings: settings}
	greedy_wins := 0
	pomcp_wins := 0
	time_termination := []int{5000}
	c_values := []float64{1}
	capture_reward := []float64{0.2}
	for _, time_limit := range time_termination {
		for _, c := range c_values {
			for _, capture_rew := range capture_reward {
				settings.Ucb_c = c
				settings.Termination_parameter = time_limit
				settings.Capture_reward = capture_rew
				iterations := 10
				for i := 0; i < iterations; i++ {
					if i == int(iterations/2) {
						temp := player1
						player1 = player2
						player2 = temp
					}
					winner, moves := server.PlayGame(&dark_pawn_chess, player1, player2)
					if winner == 1 && player1.String() == "POMCP" || winner == -1 && player2.String() == "POMCP" {
						pomcp_wins++
					} else if winner != 0 {
						greedy_wins++
					}
					// player1, player2, winner, threads, termination_condition, termination_parameter, ucb_c, moves, capture_reward
					result_string := fmt.Sprintf(
						"%v, %v, %v, %v, %v, %v, %v, %v, %v, %v\n",
						player1.String(),
						player2.String(),
						winner,
						1,
						settings.Termination_parameter,
						settings.Ucb_c,
						settings.Gamma,
						settings.Epsilon,
						moves,
						settings.Capture_reward,
					)
					// save_results(result_string)
					print(result_string)
				}
			}
		}
	}
	fmt.Printf("%v wins: %v \n", player1.String(), pomcp_wins)
	fmt.Printf("%v wins: %v \n", player2.String(), greedy_wins)
}

func save_results(result string) {
	results_file_name := "results.csv"
	f, err := os.OpenFile(results_file_name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.WriteString(result); err != nil {
		panic(err)
	}
}
