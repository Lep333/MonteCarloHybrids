package chess_variation

const row_length_lac uint = 6
const no_fields_lac uint = row_length_lac * row_length_lac
const black_base_line_start_lac uint = row_length_dpc * (row_length_dpc - 1)

type LosAlamosChess struct {
	white_pawns         uint
	white_pawns_moves   [no_fields_lac]uint
	white_pawns_capture [no_fields_lac]uint
	black_pawns         uint
	black_pawns_moves   [no_fields_lac]uint
	black_pawns_capture [no_fields_lac]uint
	white_rooks         uint
	white_rooks_moves   [no_fields_lac]uint
	black_rooks         uint
	white_knights       uint
	black_knights       uint
	knights_moves       [no_fields_lac]uint
	white_king          uint
	black_king          uint
	white_queen         uint
	black_queen         uint
	king_moves          [no_fields_lac]uint
	white_occupancy     uint
	black_occupancy     uint
	whiteToPlay         bool
	number_of_moves     int
}

func (l *LosAlamosChess) InitGame() {
	l.whiteToPlay = true
	l.number_of_moves = 0
	l.white_pawns = 0b111111000000
	l.black_pawns = 0b111111000000 << (row_length_lac * 3)
	l.white_pawns_capture = [no_fields_lac]uint{}
	l.black_pawns_capture = [no_fields_lac]uint{}
	l.white_rooks = 0b100001
	l.black_rooks = l.white_rooks << (row_length_lac * 5)
	l.white_knights = 0b010010
	l.black_knights = l.white_knights << (row_length_lac * 5)
	l.white_king = 0b000100
	l.black_king = l.white_king << (row_length_lac * 5)
	l.white_queen = 0b001000
	l.black_queen = l.white_queen << (row_length_lac * 5)
	l.knights_moves = [no_fields_lac]uint{}

	for i := uint(0); i < no_fields_lac; i++ {
		l.init_pawns(i)
		l.init_knights(i)
		l.init_kings(i)
	}

	l.set_occupancy_boards()
}

func (l *LosAlamosChess) init_pawns(i uint) {
	if i < no_fields_lac-row_length_lac {
		l.white_pawns_moves[i] = 0b1 << (i + row_length_lac)
		if i%row_length_lac != row_length_lac-1 {
			l.white_pawns_capture[i] += 0b1 << (i + row_length_lac + 1) // left capture
		}
		if i%row_length_lac != 0 {
			l.white_pawns_capture[i] += 0b1 << (i + row_length_lac - 1) // right capture
		}
	}

	if i >= row_length_lac {
		l.black_pawns_moves[i] = 0b1 << (i - row_length_lac)

		if i%row_length_lac != row_length_lac-1 {
			l.black_pawns_capture[i] += 0b1 << (i - row_length_lac + 1) // left capture
		}
		if i%row_length_lac != 0 {
			l.black_pawns_capture[i] += 0b1 << (i - row_length_lac - 1) // right capture
		}
	}
}

func (l *LosAlamosChess) init_knights(i uint) {
	col := i % row_length_lac
	row := i / row_length_lac
	// 2x up left
	if col != 0 && row < 4 {
		l.knights_moves[i] += 0b1 << (i + 2*row_length_lac - 1)
	}
	// 2x up right
	if col != 5 && row < 4 {
		l.knights_moves[i] += 0b1 << (i + 2*row_length_lac + 1)
	}
	// 2x right up
	if col < 4 && row < 5 {
		l.knights_moves[i] += 0b1 << (i + row_length_lac + 2)
	}
	// 2x right down
	if col < 4 && row > 0 {
		l.knights_moves[i] += 0b1 << (i - row_length_lac + 2)
	}
	// 2x down right
	if col < 5 && row > 1 {
		l.knights_moves[i] += 0b1 << (i - row_length_lac*2 + 1)
	}
	// 2x down left
	if col > 0 && row > 1 {
		l.knights_moves[i] += 0b1 << (i - row_length_lac*2 - 1)
	}
	// 2x left down
	if col > 1 && row > 0 {
		l.knights_moves[i] += 0b1 << (i - row_length_lac - 2)
	}
	// 2x left up
	if col > 1 && row < 5 {
		l.knights_moves[i] += 0b1 << (i + row_length_lac - 2)
	}
}

func (l *LosAlamosChess) init_kings(i uint) {
	col := i % row_length_lac
	row := i / row_length_lac
	// TODO: move possible moves in to match bit in bitboard
	// up
	if row < 5 {
		l.king_moves[i] += 0b1 << (i + row_length_lac)
	}
	// up right
	if row < 5 || col < 5 {
		l.king_moves[i] += 0b1 << (i + row_length_lac + 1)
	}
	// right
	if col < 5 {
		l.king_moves[i] += 0b1 << (i + 1)
	}
	// right down
	if row > 0 || col < 5 {
		l.king_moves[i] += 0b1 << (i - row_length_lac + 1)
	}
	// down
	if row > 0 {
		l.king_moves[i] += 0b1 << (i - row_length_lac)
	}
	// down left
	if row > 0 || col > 0 {
		l.king_moves[i] += 0b1 << (i - row_length_lac - 1)
	}
	// left
	if col > 0 {
		l.king_moves[i] += 0b1 << (i - 1)
	}
	// left up
	if row < 5 && col > 0 {
		l.king_moves[i] += 0b1 << (i + row_length_lac - 1)
	}
}

func (l *LosAlamosChess) ReturnBoard() ChessVariation {
	copy := *l
	return &copy
}

func (l *LosAlamosChess) GetPreviousBoard() ChessVariation {
	return l.ReturnBoard()
}

func (l *LosAlamosChess) GetNumberOfMoves() int {
	return l.number_of_moves
}

func (l *LosAlamosChess) PossibleMoves() []Move {
	possible_moves := []Move{}
	if l.whiteToPlay {
		for i := uint(0); i < no_fields_lac; i++ {
			// pawns
			if l.white_pawns&(0b1<<i) > 0 {
				moves_possible := (l.white_pawns_moves[i] & ^l.black_occupancy &
					^l.white_occupancy) |
					(l.white_pawns_capture[i] & l.black_occupancy)
				moves := l.move_bitboard_to_moves(i, moves_possible)
				possible_moves = append(possible_moves, moves...)
			}
			// rooks
			rook_moves := l.get_rook_moves(i)
			possible_moves = append(possible_moves, rook_moves...)
			// knights
			if l.white_knights&(0b1<<i) > 0 {
				moves_possible := l.knights_moves[i] & ^l.white_occupancy
				moves := l.move_bitboard_to_moves(i, moves_possible)
				possible_moves = append(possible_moves, moves...)
			}
			// queen
			if l.white_queen&(0b1<<i) > 0 {
				// up right
				index := i + 7
				for index < no_fields_lac && !(i%row_length_lac == 5) {
					capture := false
					if l.white_occupancy&0b1<<index > 0 {
						break
					}
					if l.black_occupancy&0b1<<index > 0 {
						capture = true
					}
					move := Move{From: int8(i), To: int8(index), Capture: capture}
					possible_moves = append(possible_moves, move)
					if index%row_length_lac == 5 {
						break
					}
				}
			}
			// king
			if l.white_king&(0b1<<i) > 0 {
				moves_possible := l.king_moves[i] & ^l.white_occupancy
				moves := l.move_bitboard_to_moves(i, moves_possible)
				possible_moves = append(possible_moves, moves...)
			}
		}
	} else {
		for i := uint(0); i < no_fields_lac; i++ {
			// pawns
			if l.black_pawns&(0b1<<i) > 0 {
				moves_possible := (l.black_pawns_moves[i] & ^l.black_occupancy &
					^l.white_occupancy) |
					(l.black_pawns_capture[i] & l.white_occupancy)
				moves := l.move_bitboard_to_moves(i, moves_possible)
				possible_moves = append(possible_moves, moves...)
			}
			// rooks
			if l.black_rooks&(0b1<<i) > 0 {
				// upwards
				index := int(i + row_length_lac)
				for index < int(no_fields_lac) {
					if l.black_occupancy&(0b1<<index) > 0 {
						break
					}
					if l.white_occupancy&(0b1<<index) > 0 {
						move := Move{From: int8(i), To: int8(index), Capture: true}
						possible_moves = append(possible_moves, move)
						break
					}
					move := Move{From: int8(i), To: int8(index), Capture: false}
					possible_moves = append(possible_moves, move)
					index += int(row_length_lac)
				}
				// down
				index = int(i) - int(row_length_lac)
				for index >= 0 {
					if l.black_occupancy&(0b1<<index) > 0 {
						break
					}
					if l.white_occupancy&(0b1<<index) > 0 {
						move := Move{From: int8(i), To: int8(index), Capture: true}
						possible_moves = append(possible_moves, move)
						break
					}
					move := Move{From: int8(i), To: int8(index), Capture: false}
					possible_moves = append(possible_moves, move)
					index += int(row_length_lac)
				}
				// right
				index = int(i) + 1
				for index < int(no_fields_lac) {
					capture := false
					if l.black_occupancy&(0b1<<index) > 0 {
						break
					}
					if l.white_occupancy&(0b1<<index) > 0 {
						capture = true
					}
					if index%int(row_length_lac) == int(row_length_lac)-1 {
						move := Move{From: int8(i), To: int8(index), Capture: capture}
						possible_moves = append(possible_moves, move)
						break
					}
					move := Move{From: int8(i), To: int8(index), Capture: false}
					possible_moves = append(possible_moves, move)
					index++
				}
				// left
				index = int(i) - 1
				for index >= 0 {
					capture := false
					if l.black_occupancy&(0b1<<index) > 0 {
						break
					}
					if l.white_occupancy&(0b1<<index) > 0 {
						capture = true
					}
					if index%int(row_length_lac) == 0 {
						move := Move{From: int8(i), To: int8(index), Capture: capture}
						possible_moves = append(possible_moves, move)
						break
					}
					move := Move{From: int8(i), To: int8(index), Capture: false}
					possible_moves = append(possible_moves, move)
					index--
				}
			}
			// knights
			if l.black_knights&(0b1<<i) > 0 {
				moves_possible := l.knights_moves[i] & ^l.black_occupancy
				moves := l.move_bitboard_to_moves(i, moves_possible)
				possible_moves = append(possible_moves, moves...)
			}
			// queen
			// TODO
			// king
			if l.black_king&(0b1<<i) > 0 {
				moves_possible := l.king_moves[i] & ^l.black_occupancy
				moves := l.move_bitboard_to_moves(i, moves_possible)
				possible_moves = append(possible_moves, moves...)
			}
		}
	}

	return possible_moves
}

func (l *LosAlamosChess) move_bitboard_to_moves(start uint, move_bitboard uint) []Move {
	possible_moves := []Move{}
	for i := uint(0); i < no_fields_lac; i++ {
		if move_bitboard&(0b1<<i) > 0 {
			capture := false
			if l.white_occupancy&l.black_occupancy&(0b1<<i) > 0 {
				capture = true
			}
			move := Move{From: int8(start), To: int8(i), Capture: capture}
			possible_moves = append(possible_moves, move)
		}
	}
	return possible_moves
}

func (l *LosAlamosChess) get_rook_moves(i uint) []Move {
	// TODO: set own_occupancy
	// set opponent_occupancy
	possible_moves := []Move{}
	if l.white_rooks&(0b1<<i) > 0 {
		// up
		index := int(i + row_length_lac)
		for index < int(no_fields_lac) {
			if l.white_occupancy&(0b1<<index) > 0 {
				break
			}
			if l.black_occupancy&(0b1<<index) > 0 {
				move := Move{From: int8(i), To: int8(index), Capture: true}
				possible_moves = append(possible_moves, move)
				break
			}
			move := Move{From: int8(i), To: int8(index), Capture: false}
			possible_moves = append(possible_moves, move)
			index += int(row_length_lac)
		}
		// down
		index = int(i) - int(row_length_lac)
		for index >= 0 {
			if l.white_occupancy&(0b1<<index) > 0 {
				break
			}
			if l.black_occupancy&(0b1<<index) > 0 {
				move := Move{From: int8(i), To: int8(index), Capture: true}
				possible_moves = append(possible_moves, move)
				break
			}
			move := Move{From: int8(i), To: int8(index), Capture: false}
			possible_moves = append(possible_moves, move)
			index += int(row_length_lac)
		}
		// right
		index = int(i) + 1
		for index < int(no_fields_lac) {
			capture := false
			if l.white_occupancy&(0b1<<index) > 0 {
				break
			}
			if l.black_occupancy&(0b1<<index) > 0 {
				capture = true
			}
			if index%int(row_length_lac) == int(row_length_lac)-1 {
				move := Move{From: int8(i), To: int8(index), Capture: capture}
				possible_moves = append(possible_moves, move)
				break
			}
			move := Move{From: int8(i), To: int8(index), Capture: false}
			possible_moves = append(possible_moves, move)
			index++
		}
		// left
		index = int(i) - 1
		for index >= 0 {
			capture := false
			if l.white_occupancy&(0b1<<index) > 0 {
				break
			}
			if l.black_occupancy&(0b1<<index) > 0 {
				capture = true
			}
			if index%int(row_length_lac) == 0 {
				move := Move{From: int8(i), To: int8(index), Capture: capture}
				possible_moves = append(possible_moves, move)
				break
			}
			move := Move{From: int8(i), To: int8(index), Capture: false}
			possible_moves = append(possible_moves, move)
			index--
		}
	}
	return possible_moves
}

func (l *LosAlamosChess) ExecuteMove(move Move) ChessVariation {
	copy := *l
	move_to_mask := uint(0b1 << move.To)
	move_from_mask := uint(0b1 << move.From)
	if l.white_rooks&move_to_mask > 0 {
		copy.white_rooks -= move_to_mask
	} else if l.white_knights&move_to_mask > 0 {
		copy.white_knights -= move_to_mask
	} else if l.white_queen&move_to_mask > 0 {
		copy.white_queen -= move_to_mask
	} else if l.white_king&move_to_mask > 0 {
		copy.white_king -= move_to_mask
	} else if l.white_pawns&move_to_mask > 0 {
		copy.white_pawns -= move_to_mask
	} else if l.black_rooks&move_to_mask > 0 {
		copy.black_rooks -= move_to_mask
	} else if l.black_knights&move_to_mask > 0 {
		copy.black_knights -= move_to_mask
	} else if l.black_queen&move_to_mask > 0 {
		copy.black_queen -= move_to_mask
	} else if l.black_king&move_to_mask > 0 {
		copy.black_king -= move_to_mask
	} else if l.black_pawns&move_to_mask > 0 {
		copy.black_pawns -= move_to_mask
	}

	if l.white_rooks&move_from_mask > 0 {
		copy.white_rooks += -move_from_mask + move_to_mask
	} else if l.white_knights&move_from_mask > 0 {
		copy.white_knights += -move_from_mask + move_to_mask
	} else if l.white_queen&move_from_mask > 0 {
		copy.white_queen += -move_from_mask + move_to_mask
	} else if l.white_king&move_from_mask > 0 {
		copy.white_king += -move_from_mask + move_to_mask
	} else if l.white_pawns&move_from_mask > 0 {
		copy.white_pawns += -move_from_mask + move_to_mask
	} else if l.black_rooks&move_from_mask > 0 {
		copy.black_rooks += -move_from_mask + move_to_mask
	} else if l.black_knights&move_from_mask > 0 {
		copy.black_knights += -move_from_mask + move_to_mask
	} else if l.black_queen&move_from_mask > 0 {
		copy.black_queen += -move_from_mask + move_to_mask
	} else if l.black_king&move_from_mask > 0 {
		copy.black_king += -move_from_mask + move_to_mask
	} else if l.black_pawns&move_from_mask > 0 {
		copy.black_pawns += -move_from_mask + move_to_mask
	}
	copy.number_of_moves++
	copy.whiteToPlay = !copy.whiteToPlay

	copy.set_occupancy_boards()
	return &copy
}

func (l *LosAlamosChess) set_occupancy_boards() {
	l.white_occupancy = l.white_rooks | l.white_knights |
		l.white_queen | l.white_king | l.white_pawns
	l.black_occupancy = l.black_rooks | l.black_knights |
		l.black_queen | l.black_king | l.black_pawns
}

func (l *LosAlamosChess) CreateView() ChessVariation {
	// TODO: implement!
	return l.ReturnBoard()
}

func (l *LosAlamosChess) GameOver() (bool, int) {
	// black wins
	if l.white_king == 0 {
		return true, -1
	}
	// white wins
	if l.black_king == 0 {
		return true, 1
	}
	return false, 0
}

func (l *LosAlamosChess) Heuristic() float64 {
	// TODO: implement!
	return 0.0
}

func (l *LosAlamosChess) String() string {
	field_string := ""
	for i := int(row_length_lac*row_length_lac) - 1; i >= 0; i-- {
		field_mask := uint(0b1 << i)
		if l.white_rooks&field_mask > 0 {
			field_string += "R "
		} else if l.white_knights&field_mask > 0 {
			field_string += "N "
		} else if l.white_queen&field_mask > 0 {
			field_string += "Q "
		} else if l.white_king&field_mask > 0 {
			field_string += "K "
		} else if l.white_pawns&field_mask > 0 {
			field_string += "P "
		} else if l.black_rooks&field_mask > 0 {
			field_string += "r "
		} else if l.black_knights&field_mask > 0 {
			field_string += "n "
		} else if l.black_queen&field_mask > 0 {
			field_string += "q "
		} else if l.black_king&field_mask > 0 {
			field_string += "k "
		} else if l.black_pawns&field_mask > 0 {
			field_string += "p "
		} else {
			field_string += "0 "
		}
		if i%int(row_length_lac) == 0 && i != 0 {
			field_string += "\n"
		}
	}
	return field_string
}
