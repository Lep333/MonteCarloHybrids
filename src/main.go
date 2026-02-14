package main

import (
	"encoding/json"
	"fmt"
	"monte_carlo_hybrids/chess_variation"
	"monte_carlo_hybrids/player"
	"monte_carlo_hybrids/server"
	"os"
	"strings"
)

var name = "DPC_IR"

func main() {
	// web_server()
	var player1, player2 player.Player
	tune_settings := player.Settings{
		Termination_parameter:     1000,
		Gamma:                     0.95,
		Epsilon:                   0.05,
		Ucb_c:                     5,
		Rollout_capture:           0.0,
		Prior_hybrid:              nil,
		Selection_hybrid:          nil,
		Rollout_selection:         nil,
		Early_playout_termination: nil,
		POMCP_name:                name,
		Opponent_modelling:        true,
		OM_Threshold:              0.4,
	}
	default_settings := player.Settings{
		Termination_parameter: 1000,
		Gamma:                 0.95,
		Epsilon:               0.05,
		Ucb_c:                 5,
		Rollout_capture:       0,
		Opponent_modelling:    true,
		OM_Threshold:          0.4,
	}
	greedy_wins := 0
	pomcp_wins := 0
	threshs := []float64{0.2, 0.4, 0.6, 0.8}
	for _, thresh := range threshs {
		tune_settings.Rollout_selection = &player.MCTS_with_informed_rollouts{
			Search_depth: 4,
			Epsilon:      thresh,
		}
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
			winner, moves, rollouts, beliefs := server.PlayGame(&game, player1, player2)
			print_game_result(player1, player2, moves, winner, rollouts, beliefs)
		}
	}
	fmt.Printf("%v wins: %v \n", player1.String(), pomcp_wins)
	fmt.Printf("%v wins: %v \n", player2.String(), greedy_wins)
}

func print_game_result(player1 player.Player, player2 player.Player,
	moves []chess_variation.Move, winner int, rollouts []int, beliefs []int) {
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
	settings1_string := DumpJSONOneLine(settings1)
	settings2_string := DumpJSONOneLine(settings2)
	settings1_string = strings.ReplaceAll(settings1_string, ",", " ")
	settings2_string = strings.ReplaceAll(settings2_string, ",", " ")
	result_string := fmt.Sprintf(
		"%v, %v, %v, %+v, %+v, %v, %v, %v \n",
		player1.String(),
		player2.String(),
		winner,
		settings1_string,
		settings2_string,
		moves,
		rollouts,
		beliefs,
	)
	print(result_string)
	save_results(result_string) // remove for saving results
}

func save_results(result string) {
	results_file_name := fmt.Sprintf("results/%v.csv", name)
	f, err := os.OpenFile(results_file_name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.WriteString(result); err != nil {
		panic(err)
	}
}

func DumpJSONOneLine(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
