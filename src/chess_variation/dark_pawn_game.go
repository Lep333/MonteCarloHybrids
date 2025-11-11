package chess_variation

import (
	"math"
	"slices"
)

const row_length uint = 5
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
	whiteToPlay         bool
	number_of_moves     int
}

func (d *DarkPawnChess) InitGame() {
	d.whiteToPlay = true
	d.number_of_moves = 0
	d.white_pawns = row_bitmask
	d.black_pawns = row_bitmask << black_base_line_start

	for i := uint(0); i < no_fields; i++ {
		if i < no_fields-row_length {
			d.white_pawns_moves[i] = 0b1 << (i + row_length)
			if i%row_length != row_length-1 {
				d.white_pawns_capture[i] += 0b1 << (i + row_length + 1) // left capture
			}
			if i%row_length != 0 {
				d.white_pawns_capture[i] += 0b1 << (i + row_length - 1) // right capture
			}
		}

		if i >= row_length {
			d.black_pawns_moves[i] = 0b1 << (i - row_length)

			if i%row_length != row_length-1 {
				d.black_pawns_capture[i] += 0b1 << (i - row_length + 1) // left capture
			}
			if i%row_length != 0 {
				d.black_pawns_capture[i] += 0b1 << (i - row_length - 1) // right capture
			}
		}
	}
}

func (d *DarkPawnChess) ReturnBoard() ChessVariation {
	copy := *d
	return &copy
}

func (d *DarkPawnChess) GetPreviousBoard() ChessVariation {
	return d.ReturnBoard()
}

func (d *DarkPawnChess) GetNumberOfMoves() int {
	return d.number_of_moves
}

// 1 white; -1 black; 0 draw
func (d *DarkPawnChess) GameOver() (bool, int) {
	game_over := false
	winner := 0

	// check for opposing pawns on base line
	for i := uint(0); i < row_length; i++ {
		if d.black_pawns&row_bitmask > 0 {
			game_over = true
			winner = -1
		}

		if d.white_pawns&(row_bitmask<<black_base_line_start) > 0 {
			game_over = true
			winner = 1
		}
	}

	if len(d.get_moves()) == 0 {
		game_over = true
	}

	return game_over, winner
}

func (d *DarkPawnChess) PossibleMoves() []Move {
	moves := []Move{}

	game_over, _ := d.GameOver()
	if game_over {
		return moves
	}

	return d.get_moves()
}

func (d *DarkPawnChess) get_moves() []Move {
	moves := []Move{}

	for i := uint(0); i < no_fields; i++ {
		if d.whiteToPlay && d.white_pawns&(0b1<<i) > 0 {
			move_to_possible := d.white_pawns_moves[i] & ^d.black_pawns
			if move_to_possible > 0 {
				move := Move{int8(i), int8(math.Log2(float64(move_to_possible))), false}
				moves = append(moves, move)
			}

			capture_possible := d.white_pawns_capture[i] & d.black_pawns
			for position := uint(0); position < no_fields; position++ {
				if (capture_possible>>position)&0b1 != 0 {
					move := Move{int8(i), int8(position), true}
					moves = append(moves, move)
				}
			}
		}

		if !d.whiteToPlay && d.black_pawns&(0b1<<i) > 0 {
			move_to_possible := d.black_pawns_moves[i] & ^d.white_pawns
			if move_to_possible > 0 {
				move := Move{int8(i), int8(math.Log2(float64(move_to_possible))), false}
				moves = append(moves, move)
			}

			capture_possible := d.black_pawns_capture[i] & d.white_pawns
			for position := uint(0); position < no_fields; position++ {
				if (capture_possible>>position)&0b1 != 0 {
					move := Move{int8(i), int8(position), true}
					moves = append(moves, move)
				}
			}
		}
	}
	return moves
}

func (d *DarkPawnChess) ExecuteMove(move Move) ChessVariation {
	copy := *d
	var mask_from uint = 1 << move.From
	var mask_to uint = 1 << move.To
	var mask uint = mask_from | mask_to
	if d.whiteToPlay {
		copy.white_pawns ^= mask
		copy.black_pawns = d.black_pawns &^ mask_to
	} else {
		copy.black_pawns ^= mask
		copy.white_pawns = d.white_pawns &^ mask_to
	}
	copy.whiteToPlay = !d.whiteToPlay
	// TODO: remove copy.prev_board = d
	copy.number_of_moves = d.number_of_moves + 1
	return &copy
}

func (d *DarkPawnChess) CreateView() ChessVariation {
	copy := *d
	poss_moves := copy.PossibleMoves()
	moves_to := []int{}
	for _, move := range poss_moves {
		if !slices.Contains(moves_to, int(move.To)) {
			moves_to = append(moves_to, int(move.To))
		}
	}

	if copy.whiteToPlay {
		copy.black_pawns = 0
		for _, move_to := range moves_to {
			copy.black_pawns += d.black_pawns & (1 << move_to)
		}
	} else {
		copy.white_pawns = 0
		for _, move_to := range moves_to {
			copy.white_pawns += d.white_pawns & (1 << move_to)
		}
	}
	return &copy
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
