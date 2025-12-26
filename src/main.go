package main

import (
	"fmt"
	"monte_carlo_hybrids/chess_variation"
	"monte_carlo_hybrids/player"
	"monte_carlo_hybrids/server"
	"os"
)

func main() {
	var player1, player2 player.Player
	tune_settings := player.Settings{
		Termination_parameter:     1000,
		Gamma:                     0.95,
		Epsilon:                   0.005,
		Ucb_c:                     5,
		Rollout_capture:           0.9,
		Rollout_selection:         &player.GreedySelection{},
		Early_playout_termination: nil,
		POMCP_name:                "POMCP-Greedy",
	}
	default_settings := player.Settings{
		Termination_parameter: 1000,
		Gamma:                 0.95,
		Epsilon:               0.005,
		Ucb_c:                 5,
		Rollout_capture:       0.9,
	}
	greedy_wins := 0
	pomcp_wins := 0
	time_termination := []int{1000}
	c_values := []float64{5}
	epsilon := []float64{0, 0.2, 0.4, 0.6, 0.8}
	for _, time_limit := range time_termination {
		for _, c := range c_values {
			for _, epsi := range epsilon {
				tune_settings.Ucb_c = c
				tune_settings.Termination_parameter = time_limit
				// tune_settings.Rollout_capture = roll_cap
				tune_settings.Rollout_selection = &player.
					GreedySelection{Epsilon: epsi}
				player1 = &player.POMCP{
					Root:            nil,
					Started_playing: false,
					Last_move:       chess_variation.Move{},
					Settings:        default_settings,
				}
				player2 = &player.POMCP{
					Root:            nil,
					Started_playing: false,
					Last_move:       chess_variation.Move{},
					Settings:        tune_settings,
				}
				iterations := 100
				for i := 0; i < iterations; i++ {
					game := chess_variation.DarkPawnChess{}
					if i == int(iterations/2) {
						temp := player1
						player1 = player2
						player2 = temp
					}
					winner, moves, rollouts := server.PlayGame(&game, player1, player2)
					print_game_result(player1, player2, moves, winner, rollouts)
				}
			}
		}
	}
	fmt.Printf("%v wins: %v \n", player1.String(), pomcp_wins)
	fmt.Printf("%v wins: %v \n", player2.String(), greedy_wins)
}

func print_game_result(player1 player.Player, player2 player.Player,
	moves []chess_variation.Move, winner int, rollouts []int) {
	// pc := reflect.ValueOf(settings1.Termination_Func).Pointer()
	// fn := runtime.FuncForPC(pc)
	// pc2 := reflect.ValueOf(settings2.Termination_Func).Pointer()
	// fn2 := runtime.FuncForPC(pc2)
	// player1, player2, winner, threads, termination_condition, termination_parameter, ucb_c, moves, capture_reward, rollout_capture
	settings1 := player.Settings{}
	settings2 := player.Settings{}
	pomcp1, ok := player1.(*player.POMCP)
	pomcp2, ok2 := player2.(*player.POMCP)
	if ok {
		settings1 = pomcp1.Settings
	}
	if ok2 {
		settings2 = pomcp2.Settings
	}
	result_string := fmt.Sprintf(
		"%v, %v, %v, %+v, %+v, %v, %v \n",
		player1.String(),
		player2.String(),
		winner,
		settings1,
		settings2,
		moves,
		rollouts,
	)
	print(result_string)
	save_results(result_string) // remove for saving results
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
