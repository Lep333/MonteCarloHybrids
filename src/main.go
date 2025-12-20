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
		Termination_parameter:  1000,
		Gamma:                  0.95,
		Epsilon:                0.005,
		Ucb_c:                  1,
		Capture_reward:         0.2,
		Rollout_capture:        0.7,
		Termination_Func:       player.MCTS_with_informed_cutoffs,
		Termination_func_param: 5.0,
	}
	default_settings := player.Settings{
		Termination_parameter: 1000,
		Gamma:                 0.95,
		Epsilon:               0.005,
		Ucb_c:                 1,
		Capture_reward:        0.0,
		Rollout_capture:       0.0,
		Termination_Func: func(
			board chess_variation.ChessVariation,
			a float64,
			b float64) (bool, float64) {
			return false, 0
		},
		Termination_func_param: 0,
	}
	greedy_wins := 0
	pomcp_wins := 0
	time_termination := []int{1000, 2000, 5000}
	c_values := []float64{1, 2, 5, 10}
	rollout_capture := []float64{0.1, 0.3, 0.5, 0.7, 0.9}
	for _, time_limit := range time_termination {
		for _, c := range c_values {
			for _, roll_cap := range rollout_capture {
				tune_settings.Ucb_c = c
				tune_settings.Termination_parameter = time_limit
				tune_settings.Rollout_capture = roll_cap
				player1 = &player.POMCP{Root: nil, Started_playing: false, Last_move: chess_variation.Move{}, Settings: default_settings}
				player2 = &player.POMCP{Root: nil, Started_playing: false, Last_move: chess_variation.Move{}, Settings: tune_settings}
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
