package chess_variation

import (
	"math"
	p "monte_carlo_hybrids/player"
)

const row_length = 5
const no_fields = row_length * row_length
const black_base_line_start = row_length * (row_length - 1)

type DarkPawnChess struct {
	white_pawns       int
	white_pawns_moves [no_fields]int
	black_pawns       int
	black_pawns_moves [no_fields]int
}

func (d *DarkPawnChess) InitGame() {
	row_bitmask := 0b11111
	d.white_pawns = row_bitmask
	d.black_pawns = row_bitmask << 20

	for i := 0; i < no_fields; i++ {
		if i < no_fields-row_length {
			d.white_pawns_moves[i] = 0b1 << (i + row_length)
		}

		if i >= row_length {
			d.black_pawns_moves[i] = 0b1 << (i - row_length)
		}
	}
}

func (d *DarkPawnChess) ReturnBoard(currPlayer p.Player) DarkPawnChess {
	return *d
}

func (d *DarkPawnChess) GameOver() bool {
	game_over := false
	row_bitmask := 0b11111

	// check for opposing pawns on base line
	for i := 0; i < row_length; i++ {
		if d.black_pawns&row_bitmask > 0 {
			game_over = true
		}

		if d.white_pawns&(row_bitmask<<black_base_line_start) > 0 {
			game_over = true
		}
	}
	return game_over
}

func (d *DarkPawnChess) PossibleMoves(whiteToPlay bool) []Move {
	moves := []Move{}
	if whiteToPlay {
		for i := 0; i < no_fields; i++ {
			if d.white_pawns&(0b1<<i) > 0 {
				move_to_possible := d.white_pawns_moves[i] & ^d.black_pawns
				move := Move{int8(i), int8(math.Log2(float64(move_to_possible)))} //
				moves = append(moves, move)
			}
		}
	} else {
		// TODO: black move generation
	}
	return moves
}
