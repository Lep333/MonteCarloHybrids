package main

import (
	"log"
	"math/rand"
	"monte_carlo_hybrids/chess_variation"
	"monte_carlo_hybrids/player"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		switch origin {
		case "http://localhost:5173",
			"http://localhost:8080":
			return true
		default:
			return true
		}
	},
}

func web_server() {
	http.HandleFunc("/play_game", play_game)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func play_game(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	tune_settings := player.Settings{
		Termination_parameter:     5000,
		Gamma:                     0.95,
		Epsilon:                   0.005,
		Ucb_c:                     1,
		Rollout_capture:           0.0,
		Prior_hybrid:              nil,
		Selection_hybrid:          nil,
		Rollout_selection:         nil,
		Early_playout_termination: nil,
		POMCP_name:                "LAC-POMCP-UCB",
	}
	var player1 player.Player
	var player2 player.Player
	var curr_player player.Player
	player1 = &player.HumanPlayer{Conn: conn}
	player2 = &player.POMCP{
		Root:            nil,
		Started_playing: false,
		Last_move:       chess_variation.Move{},
		Settings:        tune_settings,
	}
	var game chess_variation.ChessVariation
	game = &chess_variation.LosAlamosChess{}
	game.InitGame()
	if rand.Float32() > 0.5 {
		temp := player1
		player1 = player2
		player2 = temp
	}
	for {
		over, _ := game.GameOver()
		if over {
			conn.WriteJSON(player.BoardUpdate{Fen: game.FENString()})
			break
		}
		curr_player = player1
		if game.GetNumberOfMoves()%2 == 1 {
			curr_player = player2
		}
		move := curr_player.GetMove(curr_player.ViewFunc(game), game.GetNumberOfMoves()%2 == 0)
		game.ExecuteMove(move)
		log.Println(game.String())
	}
}
