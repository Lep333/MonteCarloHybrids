package player

import (
	"log"
	chess "monte_carlo_hybrids/chess_variation"

	"github.com/gorilla/websocket"
)

type HumanPlayer struct {
	Conn *websocket.Conn
}

func (h *HumanPlayer) ViewFunc(board chess.ChessVariation) chess.ChessVariation {
	white := board.GetNumberOfMoves()%2 == 0
	return board.CreateView(white)
}

type BoardUpdate struct {
	Fen   string
	Moves []chess.Move
}

func (h *HumanPlayer) GetMove(board chess.ChessVariation, whiteToPlay bool) chess.Move {
	log.Println(board.String())
	board.PossibleMoves(Moves[:])
	msg := BoardUpdate{Fen: board.FENString(), Moves: Moves[:]}
	var move chess.Move
	if err := h.Conn.WriteJSON(msg); err != nil {
		log.Println(err)
		return random_element(Moves[:])
	}
	err := h.Conn.ReadJSON(&move)
	if err != nil {
		log.Println(err)
		return random_element(Moves[:])
	}
	return move
}

func (h *HumanPlayer) String() string {
	return "Human"
}
