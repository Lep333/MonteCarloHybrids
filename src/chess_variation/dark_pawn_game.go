package chess_variation

import (
	"math"
)

const row_length uint = 8
const no_fields uint = row_length * row_length
const black_base_line_start uint = row_length * (row_length - 1)
const row_bitmask uint = (1 << row_length) - 1

type DarkPawnChess struct {
	white_pawns         uint
	white_pawns_moves   [no_fields]uint
	white_pawns_capture [no_fields]uint
	black_pawns         uint
	black_pawns_moves   [no_fields]uint
	black_pawns_capture [no_fields]uint
}

func (d *DarkPawnChess) InitGame() {
	d.white_pawns = row_bitmask
	d.black_pawns = row_bitmask << black_base_line_start

	for i := uint(0); i < no_fields; i++ {
		if i < no_fields-row_length {
			d.white_pawns_moves[i] = 0b1 << (i + row_length)
			if i-row_length+1%row_length != 0 {
				d.white_pawns_capture[i] = 0b1 << (i + row_length + 1) // left capture
			}
			if i%row_length != 0 {
				d.white_pawns_capture[i] = 0b1 << (i + row_length - 1) // right capture
			}
		}

		if i >= row_length {
			d.black_pawns_moves[i] = 0b1 << (i - row_length)

			if i-row_length+1%row_length != 0 {
				d.black_pawns_capture[i] = 0b1 << (i - row_length + 1) // left capture
			}
			if i%row_length != 0 {
				d.black_pawns_capture[i] = 0b1 << (i - row_length - 1) // right capture
			}
		}
	}
}

func (d *DarkPawnChess) ReturnBoard() ChessVariation {
	return d
}

func (d *DarkPawnChess) GameOver() bool {
	game_over := false

	// check for opposing pawns on base line
	for i := uint(0); i < row_length; i++ {
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

	for i := uint(0); i < no_fields; i++ {
		if whiteToPlay && d.white_pawns&(0b1<<i) > 0 {
			move_to_possible := d.white_pawns_moves[i] & ^d.black_pawns
			if move_to_possible > 0 {
				move := Move{int8(i), int8(math.Log2(float64(move_to_possible)))}
				moves = append(moves, move)
			}

			capture_possible := d.white_pawns_capture[i] & d.black_pawns
			for position := uint(0); position < no_fields; position++ {
				if (capture_possible>>position)&0b1 != 0 {
					move := Move{int8(i), int8(position)}
					moves = append(moves, move)
				}
			}
		}

		if !whiteToPlay && d.black_pawns&(0b1<<i) > 0 {
			move_to_possible := d.black_pawns_moves[i] & ^d.white_pawns
			if move_to_possible > 0 {
				move := Move{int8(i), int8(math.Log2(float64(move_to_possible)))}
				moves = append(moves, move)
			}

			capture_possible := d.black_pawns_capture[i] & d.white_pawns
			for position := uint(0); position < no_fields; position++ {
				if (capture_possible>>position)&0b1 != 0 {
					move := Move{int8(i), int8(position)}
					moves = append(moves, move)
				}
			}
		}
	}
	return moves
}

func (d *DarkPawnChess) ExecuteMove(whiteToPlay bool, move Move) {
	var mask_from uint = 1 << move.from
	var mask_to uint = 1 << move.to
	var mask uint = mask_from | mask_to
	if whiteToPlay {
		d.white_pawns ^= mask
		d.black_pawns = d.black_pawns &^ mask_to
	} else {
		d.black_pawns ^= mask
		d.white_pawns = d.white_pawns &^ mask_to
	}
}

func (d *DarkPawnChess) String() string {
	board := ""
	for i := int(no_fields - 1); i >= 0; i-- {
		if d.white_pawns&(0b1<<i) != 0 {
			board += "w "
		} else if d.black_pawns&(0b1<<i) != 0 {
			board += "b "
		} else {
			board += "o "
		}
		if uint(i)%row_length == 0 {
			board += "\n"
		}
	}
	return board
}
