package chess_variation

const row_length_dpc uint = 5
const no_fields_dpc uint = row_length_dpc * row_length_dpc
const black_base_line_start_dpc uint = row_length_dpc * (row_length_dpc - 1)
const row_bitmask_dpc uint = (1 << row_length_dpc) - 1

var white_pawns_moves [no_fields_dpc]uint
var white_pawns_capture [no_fields_dpc]uint
var black_pawns_moves [no_fields_dpc]uint
var black_pawns_capture [no_fields_dpc]uint

type DarkPawnChess struct {
	moves           [20]Move
	white_pawns     uint
	black_pawns     uint
	number_of_moves int
	view_mask_white uint64
	view_mask_black uint64
	move_count      int
	whiteToPlay     bool
}

func (d *DarkPawnChess) InitGame() {
	d.whiteToPlay = true
	d.number_of_moves = 0
	d.white_pawns = row_bitmask_dpc
	//d.white_pawns = 0b1111
	d.black_pawns = row_bitmask_dpc << black_base_line_start_dpc
	//d.black_pawns = 0b11000 << black_base_line_start_dpc
	white_pawns_capture = [no_fields_dpc]uint{}
	black_pawns_capture = [no_fields_dpc]uint{}

	for i := uint(0); i < no_fields_dpc; i++ {
		if i < no_fields_dpc-row_length_dpc {
			white_pawns_moves[i] = 0b1 << (i + row_length_dpc)
			if i%row_length_dpc != row_length_dpc-1 {
				white_pawns_capture[i] += 0b1 << (i + row_length_dpc + 1) // left capture
			}
			if i%row_length_dpc != 0 {
				white_pawns_capture[i] += 0b1 << (i + row_length_dpc - 1) // right capture
			}
		}

		if i >= row_length_dpc {
			black_pawns_moves[i] = 0b1 << (i - row_length_dpc)

			if i%row_length_dpc != row_length_dpc-1 {
				black_pawns_capture[i] += 0b1 << (i - row_length_dpc + 1) // left capture
			}
			if i%row_length_dpc != 0 {
				black_pawns_capture[i] += 0b1 << (i - row_length_dpc - 1) // right capture
			}
		}
	}
	d.move_count = d.get_moves(d.moves[:])
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
	for i := uint(0); i < row_length_dpc; i++ {
		if d.black_pawns&row_bitmask_dpc > 0 {
			game_over = true
			winner = -1
		}

		if d.white_pawns&(row_bitmask_dpc<<black_base_line_start_dpc) > 0 {
			game_over = true
			winner = 1
		}
	}

	if d.get_moves(d.moves[:]) == 0 {
		game_over = true
		no_of_white_pieces := 0
		no_of_black_pieces := 0
		for i := uint(0); i < no_fields_dpc; i++ {
			if d.white_pawns>>i&1 == 1 {
				no_of_white_pieces++
			}
			if d.black_pawns>>i&1 == 1 {
				no_of_black_pieces++
			}
		}
		if no_of_white_pieces > no_of_black_pieces {
			winner = 1
		} else if no_of_black_pieces > no_of_white_pieces {
			winner = -1
		}
	}

	return game_over, winner
}

func (d *DarkPawnChess) PossibleMoves() []Move {
	game_over, _ := d.GameOver()
	if game_over {
		d.move_count = 0
		return d.moves[:d.move_count]
	}
	move_count := d.get_moves(d.moves[:])
	return d.moves[:move_count]
}

func (d *DarkPawnChess) get_moves(moves []Move) int {
	n := 0
	for i := uint(0); i < no_fields_dpc; i++ {
		if d.whiteToPlay && d.white_pawns&(0b1<<i) > 0 {
			move_to_possible := white_pawns_moves[i] & ^d.black_pawns & ^d.white_pawns
			if move_to_possible > 0 {
				move := Move{int8(i), int8(i + row_length_dpc), false}
				moves[n] = move
				n++
			}

			capture_possible := white_pawns_capture[i] & d.black_pawns
			for position := uint(0); position < no_fields_dpc; position++ {
				if (capture_possible>>position)&0b1 != 0 {
					move := Move{int8(i), int8(position), true}
					moves[n] = move
					n++
				}
			}
		}

		if !d.whiteToPlay && d.black_pawns&(0b1<<i) > 0 {
			move_to_possible := black_pawns_moves[i] & ^d.white_pawns & ^d.black_pawns
			if move_to_possible > 0 {
				move := Move{int8(i), int8(i - row_length_dpc), false}
				moves[n] = move
				n++
			}

			capture_possible := black_pawns_capture[i] & d.white_pawns
			for position := uint(0); position < no_fields_dpc; position++ {
				if (capture_possible>>position)&0b1 != 0 {
					move := Move{int8(i), int8(position), true}
					moves[n] = move
					n++
				}
			}
		}
	}
	return n
}

func (d *DarkPawnChess) set_vision() {
	white_mask := uint(0)
	black_mask := uint(0)
	for i := uint(0); i < no_fields_dpc; i++ {
		if d.white_pawns&(0b1<<i) > 0 {
			white_mask |= white_pawns_moves[i]
			white_mask |= white_pawns_capture[i]
		}

		if d.black_pawns&(0b1<<i) > 0 {
			black_mask |= black_pawns_moves[i]
			black_mask |= black_pawns_capture[i]
		}
	}
	d.view_mask_white = uint64(white_mask)
	d.view_mask_black = uint64(black_mask)
}

// executes move and updates view_mask
func (d *DarkPawnChess) ExecuteMove(move Move) {
	var mask_from uint = 1 << move.From
	var mask_to uint = 1 << move.To
	var mask uint = mask_from | mask_to
	if d.whiteToPlay {
		d.white_pawns ^= mask
		d.black_pawns = d.black_pawns &^ mask_to
	} else {
		d.white_pawns = d.white_pawns &^ mask_to
		d.black_pawns ^= mask
	}
	d.whiteToPlay = !d.whiteToPlay
	d.number_of_moves = d.number_of_moves + 1
	d.move_count = d.get_moves(d.moves[:])
	d.create_view_mask()
}

func (d *DarkPawnChess) UndoMove(move Move) {
	d.whiteToPlay = !d.whiteToPlay
	d.number_of_moves = d.number_of_moves - 1
	var mask_from uint = 1 << move.From
	var mask_to uint = 1 << move.To
	var mask uint = mask_from | mask_to
	if d.whiteToPlay {
		d.white_pawns ^= mask
		if move.Capture {
			d.black_pawns = d.black_pawns | mask_to
		}
	} else {
		if move.Capture {
			d.white_pawns = d.white_pawns | mask_to
		}
		d.black_pawns ^= mask
	}
	d.move_count = d.get_moves(d.moves[:])
	d.create_view_mask()
}

var static_moves = [20]Move{}

func (d *DarkPawnChess) CreateView(white bool) ChessVariation {
	copy := *d
	copy.set_vision()

	if white {
		copy.black_pawns = copy.black_pawns & uint(d.view_mask_white)
	} else {
		copy.white_pawns = copy.white_pawns & uint(d.view_mask_black)
	}
	return &copy
}

func (d *DarkPawnChess) create_view_mask() {
	d.set_vision()
}

func (d *DarkPawnChess) ViewHash(white bool) uint64 {
	hash := uint64(0)

	if white {
		hash = (uint64(d.black_pawns) & d.view_mask_white) << 32
		hash |= uint64(d.white_pawns)
	} else {
		hash = uint64(d.white_pawns) & d.view_mask_black
		hash |= uint64(d.black_pawns) << 32
	}
	return hash
}

func (d *DarkPawnChess) Hash() uint64 {
	return (uint64(d.black_pawns) << 32) | uint64(d.white_pawns)
}

func (d *DarkPawnChess) Heuristic(white bool) float64 {
	if over, val := d.GameOver(); over {
		if white {
			return float64(val)
		} else {
			return -float64(val)
		}
	}
	value := 0.0
	no_white_pawns := 0
	no_black_pawns := 0
	white_coverage := 0
	black_coverage := 0
	for i := uint(0); i < no_fields_dpc; i++ {
		col := i % row_length_dpc
		row := i / row_length_dpc
		if (d.black_pawns>>i)&0b1 == 1 {
			no_black_pawns++
			// up right
			if row < row_length_dpc-1 && col < row_length_dpc-1 && d.black_pawns&(1<<i+row_length_dpc+1) > 0 {
				black_coverage++
			}
			// up left
			if row < row_length_dpc-1 && col > 0 && d.black_pawns&(1<<i+row_length_dpc-1) > 0 {
				black_coverage++
			}
		}
		if (d.white_pawns>>i)&0b1 == 1 {
			no_white_pawns++
			// down right
			if row > 0 && col < row_length_dpc-1 && d.white_pawns&(1<<i-row_length_dpc+1) > 0 {
				white_coverage++
			}
			// down left
			if row > 0 && col > 0 && d.white_pawns&(1<<i-row_length_dpc-1) > 0 {
				white_coverage++
			}
		}
	}

	var material_advantage float64
	var coverage float64
	if white {
		material_advantage = float64(no_white_pawns-no_black_pawns) / float64(row_length_dpc)
		if white_coverage+black_coverage == 0 {
			coverage = 0
		} else {
			coverage = float64(white_coverage-black_coverage) / float64(row_length_dpc)
		}
	} else {
		material_advantage = float64(no_black_pawns-no_white_pawns) / float64(row_length_dpc)
		if white_coverage+black_coverage == 0 {
			coverage = 0
		} else {
			coverage = float64(black_coverage-white_coverage) / float64(row_length_dpc)
		}
	}
	value = (material_advantage + coverage) / 2
	return value
}

func (d *DarkPawnChess) String() string {
	board := ""
	for i := int(no_fields_dpc - 1); i >= 0; i-- {
		if d.white_pawns&(0b1<<i) != 0 {
			board += "w "
		} else if d.black_pawns&(0b1<<i) != 0 {
			board += "b "
		} else {
			board += "o "
		}
		if uint(i)%row_length_dpc == 0 {
			board += "\n"
		}
	}
	return board
}

func (d *DarkPawnChess) FENString() string {
	board := ""
	for i := int(no_fields_dpc - 1); i >= 0; i-- {
		if d.white_pawns&(0b1<<i) != 0 {
			board += "w "
		} else if d.black_pawns&(0b1<<i) != 0 {
			board += "b "
		} else {
			board += "o "
		}
		if uint(i)%row_length_dpc == 0 {
			board += "\n"
		}
	}
	return board
}

func (d *DarkPawnChess) Set_Board(white_pawns uint, black_pawns uint, number_of_moves int, whiteToplay bool) {
	d.white_pawns = white_pawns
	d.black_pawns = black_pawns
	d.number_of_moves = number_of_moves
	d.whiteToPlay = true
}
