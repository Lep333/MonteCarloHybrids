package chess_variation

import (
	"math"
	"math/bits"
	"math/rand/v2"
	"strconv"
)

const row_length_dc uint64 = 8
const no_fields_dc uint64 = row_length_dc * row_length_dc
const black_base_line_start_dc uint64 = row_length_dc * (row_length_dc - 1)
const no_of_piece_types_dc = 6
const len_zobrist_numbers_dc = no_fields_dc * no_of_piece_types_dc * 2

var dc_white_pawns_moves [no_fields_dc]uint64
var dc_white_pawns_capture [no_fields_dc]uint64
var dc_black_pawns_moves [no_fields_dc]uint64
var dc_black_pawns_capture [no_fields_dc]uint64
var dc_knights_moves [no_fields_dc]uint64
var dc_king_moves [no_fields_dc]uint64
var dc_zobrist_numbers [len_zobrist_numbers_dc]uint64

type DarkChess struct {
	white_pawns     uint64
	black_pawns     uint64
	white_rooks     uint64
	black_rooks     uint64
	white_knights   uint64
	black_knights   uint64
	white_king      uint64
	black_king      uint64
	white_queen     uint64
	black_queen     uint64
	white_bishop    uint64
	black_bishop    uint64
	white_occupancy uint64
	black_occupancy uint64
	whiteToPlay     bool
	number_of_moves int
	move_count      int
}

func init() {
	for i := 0; i < int(len_zobrist_numbers_dc); i++ {
		dc_zobrist_numbers[i] = rand.Uint64()
	}
	for i := uint64(0); i < no_fields_dc; i++ {
		dc_init_pawns(i)
		dc_init_knights(i)
		dc_init_kings(i)
	}
}

func (dc *DarkChess) InitGame() {
	dc.whiteToPlay = true
	dc.number_of_moves = 0
	dc.white_pawns = uint64(0b11111111) << row_length_dc
	dc.black_pawns = uint64(0b11111111) << (row_length_dc * 6)
	dc.white_rooks = uint64(0b10000001)
	dc.black_rooks = uint64(dc.white_rooks) << (row_length_dc * 7)
	dc.white_knights = uint64(0b01000010)
	dc.black_knights = uint64(dc.white_knights) << (row_length_dc * 7)
	dc.white_king = uint64(0b10000)
	dc.black_king = uint64(dc.white_king) << (row_length_dc * 7)
	dc.white_queen = uint64(0b01000)
	dc.black_queen = uint64(dc.white_queen) << (row_length_dc * 7)
	dc.white_bishop = uint64(0b00100100)
	dc.black_bishop = uint64(dc.white_bishop) << (row_length_dc * 7)
	dc.set_occupancy_boards()
}

func dc_init_pawns(i uint64) {
	if i < no_fields_dc-row_length_dc {
		dc_white_pawns_moves[i] = 0b1 << (i + row_length_dc)
		if i%row_length_dc != row_length_dc-1 {
			dc_white_pawns_capture[i] += 0b1 << (i + row_length_dc + 1) // left capture
		}
		if i%row_length_dc != 0 {
			dc_white_pawns_capture[i] += 0b1 << (i + row_length_dc - 1) // right capture
		}
	}

	if i >= row_length_dc {
		dc_black_pawns_moves[i] = 0b1 << (i - row_length_dc)

		if i%row_length_dc != row_length_dc-1 {
			dc_black_pawns_capture[i] += 0b1 << (i - row_length_dc + 1) // left capture
		}
		if i%row_length_dc != 0 {
			dc_black_pawns_capture[i] += 0b1 << (i - row_length_dc - 1) // right capture
		}
	}
}

func dc_init_knights(i uint64) {
	col := i % row_length_dc
	row := i / row_length_dc
	// TODO:
	// 2x up left
	if col != 0 && row < row_length_dc-2 {
		dc_knights_moves[i] += 0b1 << (i + 2*row_length_dc - 1)
	}
	// 2x up right
	if col != row_length_dc-1 && row < row_length_dc-2 {
		dc_knights_moves[i] += 0b1 << (i + 2*row_length_dc + 1)
	}
	// 2x right up
	if col < row_length_dc-2 && row < row_length_dc-1 {
		dc_knights_moves[i] += 0b1 << (i + row_length_dc + 2)
	}
	// 2x right down
	if col < row_length_dc-2 && row > 0 {
		dc_knights_moves[i] += 0b1 << (i - row_length_dc + 2)
	}
	// 2x down right
	if col < row_length_dc-1 && row > 1 {
		dc_knights_moves[i] += 0b1 << (i - row_length_dc*2 + 1)
	}
	// 2x down left
	if col > 0 && row > 1 {
		dc_knights_moves[i] += 0b1 << (i - row_length_dc*2 - 1)
	}
	// 2x left down
	if col > 1 && row > 0 {
		dc_knights_moves[i] += 0b1 << (i - row_length_dc - 2)
	}
	// 2x left up
	if col > 1 && row < row_length_dc-1 {
		dc_knights_moves[i] += 0b1 << (i + row_length_dc - 2)
	}
}

func dc_init_kings(i uint64) {
	col := i % row_length_dc
	row := i / row_length_dc

	// up
	if row < row_length_dc-1 {
		dc_king_moves[i] += 0b1 << (i + row_length_dc)
	}
	// up right
	if row < row_length_dc-1 && col < row_length_dc-1 {
		dc_king_moves[i] += 0b1 << (i + row_length_dc + 1)
	}
	// right
	if col < row_length_dc-1 {
		dc_king_moves[i] += 0b1 << (i + 1)
	}
	// right down
	if row > 0 && col < row_length_dc-1 {
		dc_king_moves[i] += 0b1 << (i - row_length_dc + 1)
	}
	// down
	if row > 0 {
		dc_king_moves[i] += 0b1 << (i - row_length_dc)
	}
	// down left
	if row > 0 && col > 0 {
		dc_king_moves[i] += 0b1 << (i - row_length_dc - 1)
	}
	// left
	if col > 0 {
		dc_king_moves[i] += 0b1 << (i - 1)
	}
	// left up
	if row < row_length_dc-1 && col > 0 {
		dc_king_moves[i] += 0b1 << (i + row_length_dc - 1)
	}
}

func (l *DarkChess) ReturnBoard() ChessVariation {
	copy := *l
	return &copy
}

func (l *DarkChess) GetPreviousBoard() ChessVariation {
	return l.ReturnBoard()
}

func (l *DarkChess) GetNumberOfMoves() int {
	return l.number_of_moves
}

func (l *DarkChess) PossibleMoves() []Move {
	moves := [400]Move{}
	l.set_occupancy_boards()
	l.move_count = l.generate_moves(moves[:])
	return moves[:l.move_count]
}

func (dc *DarkChess) generate_moves(buffer []Move) int {
	n := 0
	own_pawns := dc.white_pawns
	own_pawns_moves := dc_white_pawns_moves
	own_pawns_capture := dc_white_pawns_capture
	own_rooks := dc.white_rooks
	own_knights := dc.white_knights
	own_queen := dc.white_queen
	own_king := dc.white_king
	own_bishops := dc.white_bishop
	own_occupancy := dc.white_occupancy
	opponent_occupancy := dc.black_occupancy
	white_to_play := true
	if !dc.whiteToPlay {
		own_pawns = dc.black_pawns
		own_pawns_moves = dc_black_pawns_moves
		own_pawns_capture = dc_black_pawns_capture
		own_rooks = dc.black_rooks
		own_knights = dc.black_knights
		own_queen = dc.black_queen
		own_king = dc.black_king
		own_bishops = dc.black_bishop
		own_occupancy = dc.black_occupancy
		opponent_occupancy = dc.white_occupancy
		white_to_play = false
	}

	for i := uint64(0); i < no_fields_dc; i++ {
		// pawns
		if own_pawns&(0b1<<i) > 0 {
			moves_possible := (own_pawns_moves[i] & ^opponent_occupancy &
				^own_occupancy) |
				(own_pawns_capture[i] & opponent_occupancy)
			n = dc.move_bitboard_to_moves(i, moves_possible, buffer, n, opponent_occupancy)
			// pawns to in front
			if dc.whiteToPlay {
				if i < 16 && (1<<(i+8)&(dc.white_occupancy|dc.black_occupancy)) == 0 &&
					(uint64(1)<<(i+16)&(dc.white_occupancy|dc.black_occupancy)) == 0 {
					n = dc.move_bitboard_to_moves(i, 1<<(i+16), buffer, n, opponent_occupancy)
				}
			} else {
				if i > 47 && (1<<(i-8)&(dc.white_occupancy|dc.black_occupancy)) == 0 &&
					(uint64(1)<<(i-16)&(dc.white_occupancy|dc.black_occupancy)) == 0 {
					n = dc.move_bitboard_to_moves(i, 1<<(i-16), buffer, n, opponent_occupancy)
				}
			}
		}
		// rooks
		if own_rooks&(0b1<<i) > 0 {
			n = dc.get_rook_moves(i, white_to_play, buffer, n)
		}
		// knights
		if own_knights&(0b1<<i) > 0 {
			moves_possible := dc_knights_moves[i] & ^own_occupancy
			n = dc.move_bitboard_to_moves(i, moves_possible, buffer, n, opponent_occupancy)
		}
		// queen
		if own_queen&(0b1<<i) > 0 {
			n = dc.get_rook_moves(i, white_to_play, buffer, n)
			n = dc.get_bishop_moves(i, white_to_play, buffer, n)
		}
		// king
		if own_king&(0b1<<i) > 0 {
			moves_possible := dc_king_moves[i] & (^own_occupancy)
			n = dc.move_bitboard_to_moves(i, moves_possible, buffer, n, opponent_occupancy)
		}
		// bishops
		if own_bishops&(0b1<<i) > 0 {
			n = dc.get_bishop_moves(i, white_to_play, buffer, n)
		}
	}
	return n
}

func (l *DarkChess) move_bitboard_to_moves(start uint64, move_bitboard uint64, buffer []Move, n int, opp_occupancy uint64) int {
	for i := uint64(0); i < no_fields_dc; i++ {
		move_to := uint64(0b1 << i)
		if move_bitboard&move_to > 0 {
			capture := false
			if opp_occupancy&move_to > 0 {
				capture = true
			}
			move := Move{From: int8(start), To: int8(i), Capture: capture}
			buffer[n] = move
			n++
		}
	}
	return n
}

func (l *DarkChess) get_captured_piece(move_to uint64) Piece {
	var piece_type Piece
	opp_king := l.black_king
	opp_queen := l.black_queen
	opp_pawn := l.black_pawns
	opp_rook := l.black_rooks
	opp_knights := l.black_knights
	opp_bishops := l.black_bishop
	if !l.whiteToPlay {
		opp_king = l.white_king
		opp_queen = l.white_queen
		opp_pawn = l.white_pawns
		opp_rook = l.white_rooks
		opp_knights = l.white_knights
		opp_bishops = l.white_bishop
	}
	if opp_king&move_to > 0 {
		piece_type = King
	} else if opp_queen&move_to > 0 {
		piece_type = Queen
	} else if opp_pawn&move_to > 0 {
		piece_type = Pawn
	} else if opp_rook&move_to > 0 {
		piece_type = Rook
	} else if opp_knights&move_to > 0 {
		piece_type = Knight
	} else if opp_bishops&move_to > 0 {
		piece_type = Bishop
	}
	return piece_type
}

func (l *DarkChess) get_rook_moves(i uint64, white_to_play bool, buffer []Move, n int) int {
	own_occupancy := l.white_occupancy
	opponent_occupancy := l.black_occupancy
	if !white_to_play {
		own_occupancy = l.black_occupancy
		opponent_occupancy = l.white_occupancy
	}
	// up
	index := int(i + row_length_dc)
	for index < int(no_fields_dc) && !(i/row_length_dc == row_length_dc-1) {
		move_to := uint64(0b1 << index)
		if own_occupancy&(move_to) > 0 {
			break
		}
		if opponent_occupancy&(move_to) > 0 {
			move := Move{From: int8(i), To: int8(index), Capture: true}
			buffer[n] = move
			n++
			break
		}
		move := Move{From: int8(i), To: int8(index), Capture: false}
		buffer[n] = move
		n++
		index += int(row_length_dc)
	}
	// down
	index = int(i) - int(row_length_dc)
	for index >= 0 && !(i/row_length_dc == 0) {
		move_to := uint64(0b1 << index)
		if own_occupancy&(move_to) > 0 {
			break
		}
		if opponent_occupancy&(0b1<<index) > 0 {
			move := Move{From: int8(i), To: int8(index), Capture: true}
			buffer[n] = move
			n++
			break
		}
		move := Move{From: int8(i), To: int8(index), Capture: false}
		buffer[n] = move
		n++
		index -= int(row_length_dc)
	}
	// right
	index = int(i) + 1
	for index < int(no_fields_dc) && !(i%row_length_dc == row_length_dc-1) {
		capture := false
		move_to := uint64(0b1 << index)
		if own_occupancy&(move_to) > 0 {
			break
		}
		if opponent_occupancy&(move_to) > 0 {
			capture = true
		}
		if index%int(row_length_dc) == int(row_length_dc)-1 || capture {
			move := Move{From: int8(i), To: int8(index), Capture: capture}
			buffer[n] = move
			n++
			break
		}
		move := Move{From: int8(i), To: int8(index), Capture: false}
		buffer[n] = move
		n++
		index++
	}
	// left
	index = int(i) - 1
	for index >= 0 && !(i%row_length_dc == 0) {
		capture := false
		move_to := uint64(0b1 << index)
		if own_occupancy&(move_to) > 0 {
			break
		}
		if opponent_occupancy&(move_to) > 0 {
			capture = true
		}
		if index%int(row_length_dc) == 0 || capture {
			move := Move{From: int8(i), To: int8(index), Capture: capture}
			buffer[n] = move
			n++
			break
		}
		move := Move{From: int8(i), To: int8(index), Capture: false}
		buffer[n] = move
		n++
		index--
	}
	return n
}

func (l *DarkChess) get_bishop_moves(i uint64, white_to_play bool, buffer []Move, n int) int {
	own_occupancy := l.white_occupancy
	opponent_occupancy := l.black_occupancy
	if !l.whiteToPlay {
		own_occupancy = l.black_occupancy
		opponent_occupancy = l.white_occupancy
	}
	offsets := []int{9, 7, -7, -9}
	last_row_start := uint64(no_fields_dc - row_length_dc)
	for _, offset := range offsets {
		index := i
		capture := false
		for {
			col := int(index) % int(row_length_dc)
			//row := int(index) / int(row_length_dc)
			if offset == -9 && (index < row_length_dc || col == 0) {
				break
			}
			if offset == -7 && (index < row_length_dc || col == int(row_length_dc)-1) {
				break
			}
			if offset == 7 && (index >= last_row_start || col == 0) {
				break
			}
			if offset == 9 && (index >= last_row_start || col == int(row_length_dc)-1) {
				break
			}
			index += uint64(offset)
			move_to := 1 << index
			if own_occupancy&uint64(move_to) > 0 {
				break
			}
			capture = opponent_occupancy&uint64(move_to) > 0
			move := Move{From: int8(i), To: int8(index), Capture: capture}
			buffer[n] = move
			n++

			if capture {
				break
			}
		}
	}
	return n
}

func (l *DarkChess) compute_vision(white bool) uint64 {
	flip_players_turn := false
	local_moves := [400]Move{}
	vision := uint64(0)
	if white != l.whiteToPlay {
		l.whiteToPlay = !l.whiteToPlay
		flip_players_turn = true
	}
	move_count := l.generate_moves(local_moves[:])
	for _, move := range local_moves[:move_count] {
		vision |= (1 << move.To) | (1 << move.From)
	}
	if white {
		field_in_front_of_pawns := l.white_pawns << row_length_dc
		vision |= l.white_pawns | l.white_rooks | l.white_knights |
			l.white_queen | l.white_king | field_in_front_of_pawns | l.white_bishop
	} else {
		field_in_front_of_pawns := l.black_pawns >> row_length_dc
		vision |= l.black_pawns | l.black_rooks | l.black_knights |
			l.black_queen | l.black_king | field_in_front_of_pawns | l.black_bishop
	}
	if flip_players_turn {
		l.whiteToPlay = !l.whiteToPlay
	}
	return vision
}

func (l *DarkChess) ExecuteMove(move Move) {
	last_row_start := 56
	first_row_end := 7
	move_to_mask := uint64(0b1 << move.To)
	move_from_mask := uint64(0b1 << move.From)
	l.white_rooks = l.white_rooks &^ move_to_mask
	l.white_knights = l.white_knights &^ move_to_mask
	l.white_queen = l.white_queen &^ move_to_mask
	l.white_king = l.white_king &^ move_to_mask
	l.white_pawns = l.white_pawns &^ move_to_mask
	l.white_bishop = l.white_bishop &^ move_to_mask
	l.black_rooks = l.black_rooks &^ move_to_mask
	l.black_knights = l.black_knights &^ move_to_mask
	l.black_queen = l.black_queen &^ move_to_mask
	l.black_king = l.black_king &^ move_to_mask
	l.black_pawns = l.black_pawns &^ move_to_mask
	l.black_bishop = l.black_bishop &^ move_to_mask

	if l.white_rooks&move_from_mask > 0 {
		l.white_rooks = (l.white_rooks &^ move_from_mask) | move_to_mask
	}
	if l.white_knights&move_from_mask > 0 {
		l.white_knights = (l.white_knights &^ move_from_mask) | move_to_mask
	}
	if l.white_queen&move_from_mask > 0 {
		l.white_queen = (l.white_queen &^ move_from_mask) | move_to_mask
	}
	if l.white_king&move_from_mask > 0 {
		l.white_king = (l.white_king &^ move_from_mask) | move_to_mask
	}
	if l.white_pawns&move_from_mask > 0 {
		if move.To >= int8(last_row_start) {
			l.white_queen = l.white_queen | move_to_mask
			l.white_pawns = l.white_pawns &^ move_from_mask
		} else {
			l.white_pawns = (l.white_pawns &^ move_from_mask) | move_to_mask
		}
	}
	if l.white_bishop&move_from_mask > 0 {
		l.white_bishop = (l.white_bishop &^ move_from_mask) | move_to_mask
	}
	if l.black_rooks&move_from_mask > 0 {
		l.black_rooks = (l.black_rooks &^ move_from_mask) | move_to_mask
	}
	if l.black_knights&move_from_mask > 0 {
		l.black_knights = (l.black_knights &^ move_from_mask) | move_to_mask
	}
	if l.black_queen&move_from_mask > 0 {
		l.black_queen = (l.black_queen &^ move_from_mask) | move_to_mask
	}
	if l.black_king&move_from_mask > 0 {
		l.black_king = (l.black_king &^ move_from_mask) | move_to_mask
	}
	if l.black_pawns&move_from_mask > 0 {
		if move.To <= int8(first_row_end) {
			l.black_queen = l.black_queen | move_to_mask
			l.black_pawns = l.black_pawns &^ move_from_mask
		} else {
			l.black_pawns = (l.black_pawns &^ move_from_mask) | move_to_mask
		}
	}
	if l.black_bishop&move_from_mask > 0 {
		l.black_bishop = (l.black_bishop &^ move_from_mask) | move_to_mask
	}
	l.number_of_moves++
	l.whiteToPlay = !l.whiteToPlay

	l.set_occupancy_boards()
}

func (dc *DarkChess) set_occupancy_boards() {
	dc.white_occupancy = dc.white_rooks | dc.white_knights |
		dc.white_queen | dc.white_king | dc.white_pawns | dc.white_bishop
	dc.black_occupancy = dc.black_rooks | dc.black_knights |
		dc.black_queen | dc.black_king | dc.black_pawns | dc.black_bishop
}

func (l *DarkChess) CreateView(white bool) ChessVariation {
	copy := *l
	moves := l.PossibleMoves()
	vision := uint64(0)
	for _, move := range moves {
		vision |= (1 << move.To)
	}
	if white {
		field_in_front_of_pawns := l.white_pawns << row_length_dc
		vision |= l.white_pawns | l.white_rooks | l.white_knights |
			l.white_queen | l.white_king | field_in_front_of_pawns | l.white_bishop
	} else {
		field_in_front_of_pawns := l.black_pawns >> row_length_dc
		vision |= l.black_pawns | l.black_rooks | l.black_knights |
			l.black_queen | l.black_king | field_in_front_of_pawns | l.black_bishop
	}
	copy.white_pawns &= vision
	copy.white_rooks &= vision
	copy.white_knights &= vision
	copy.white_queen &= vision
	copy.white_king &= vision
	copy.white_bishop &= vision
	copy.black_pawns &= vision
	copy.black_rooks &= vision
	copy.black_knights &= vision
	copy.black_queen &= vision
	copy.black_king &= vision
	copy.black_bishop &= vision
	copy.set_occupancy_boards()

	return &copy
}

func (l *DarkChess) GetView(white bool) uint64 {
	return uint64(l.compute_vision(white))
}

func (l *DarkChess) Create_fallback_particle(belief ChessVariation, white bool) ChessVariation {
	if concrete_belief, ok := belief.(*DarkChess); ok {
		copy := *l
		view := l.GetView(white)
		if white {
			no_information := concrete_belief.black_occupancy & ^view
			//copy.black_bishop |= no_information & concrete_belief.black_bishop
			copy.black_king |= no_information & concrete_belief.black_king
			copy.black_knights |= no_information & concrete_belief.black_knights
			copy.black_pawns |= no_information & concrete_belief.black_pawns
			copy.black_queen |= no_information & concrete_belief.black_queen
			copy.black_rooks |= no_information & concrete_belief.black_rooks
			// iterate over belief and copy and check if they are equal (after last move is executed in belief)
		} else {
			no_information := concrete_belief.white_occupancy & ^view
			//copy.white_bishop |= no_information & concrete_belief.white_bishop
			copy.white_king |= no_information & concrete_belief.white_king
			copy.white_knights |= no_information & concrete_belief.white_knights
			copy.white_pawns |= no_information & concrete_belief.white_pawns
			copy.white_queen |= no_information & concrete_belief.white_queen
			copy.white_rooks |= no_information & concrete_belief.white_rooks
		}
		copy.set_occupancy_boards()
		return &copy
	}
	panic("Didnt call Create fallback particle with a belief of type Dark Chess!")
}

func (l *DarkChess) ViewHash(white bool) uint64 {
	hash := uint64(0)
	mask := l.compute_vision(white)
	gap := no_of_piece_types * 2

	for i := 0; i < int(no_fields_dc); i++ {
		if l.white_pawns&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap]
		} else if l.white_rooks&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+int(Rook)]
		} else if l.white_knights&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+int(Knight)]
		} else if l.white_queen&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+int(Queen)]
		} else if l.white_king&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+int(King)]
		} else if l.white_bishop&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+int(Bishop)]
		} else if l.black_pawns&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types_dc]
		} else if l.black_rooks&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types_dc+int(Rook)]
		} else if l.black_knights&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types_dc+int(Knight)]
		} else if l.black_queen&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types_dc+int(Queen)]
		} else if l.black_king&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types_dc+int(King)]
		} else if l.black_bishop&mask&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types_dc+int(Bishop)]
		}
	}
	return hash
}

func (l *DarkChess) Hash() uint64 {
	hash := uint64(0)
	gap := no_of_piece_types * 2

	for i := 0; i < int(no_fields_dc); i++ {
		if l.white_pawns&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap]
		} else if l.white_rooks&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+int(Rook)]
		} else if l.white_knights&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+int(Knight)]
		} else if l.white_queen&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+int(Queen)]
		} else if l.white_king&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+int(King)]
		} else if l.white_bishop&(1<<i) > 0 {
			//hash ^= dc_zobrist_numbers[i*gap+int(Bishop)]
		} else if l.black_pawns&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types]
		} else if l.black_rooks&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types+int(Rook)]
		} else if l.black_knights&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types+int(Knight)]
		} else if l.black_queen&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types+int(Queen)]
		} else if l.black_king&(1<<i) > 0 {
			hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types+int(King)]
		} else if l.black_bishop&(1<<i) > 0 {
			//hash ^= dc_zobrist_numbers[i*gap+no_of_piece_types+int(Bishop)]
		}
	}
	return hash
}

func (l *DarkChess) GameOver() (bool, int) {
	// black wins
	if l.white_king == 0 {
		return true, -1
	}
	// white wins
	if l.black_king == 0 {
		return true, 1
	}

	max_moves := 200
	if l.number_of_moves > max_moves {
		return true, 0
	}

	return false, 0
}

func (l *DarkChess) Heuristic(white bool) float64 {
	over, result := l.GameOver()
	if over {
		if !white {
			result = -result
		}
		return float64(result)
	}
	value := 0.0
	white_material := 0.0
	black_material := 0.0
	for i := 0; i < int(no_fields_dc); i++ {
		if l.white_queen&(0b1<<i) > 0 {
			white_material += 9
		} else if l.white_rooks&(0b1<<i) > 0 {
			white_material += 5
		} else if l.white_knights&(0b1<<i) > 0 {
			white_material += 3
		} else if l.white_bishop&(0b1<<i) > 0 {
			white_material += 3
		} else if l.white_pawns&(0b1<<i) > 0 {
			white_material += 1
		} else if l.black_queen&(0b1<<i) > 0 {
			black_material += 9
		} else if l.black_rooks&(0b1<<i) > 0 {
			black_material += 5
		} else if l.black_knights&(0b1<<i) > 0 {
			black_material += 3
		} else if l.black_bishop&(0b1<<i) > 0 {
			black_material += 3
		} else if l.black_pawns&(0b1<<i) > 0 {
			black_material += 1
		}
	}
	white_mobility := 0
	black_mobility := 0
	white_mobility = bits.OnesCount(uint(l.compute_vision(true)))
	black_mobility = bits.OnesCount(uint(l.compute_vision(false)))
	avg_material := (white_material + black_material) / 2
	value = (white_material - black_material + float64(l.pawn_structure()) + 0.1*float64(white_mobility-black_mobility)) / avg_material
	value = math.Max(math.Min(value, 1), -1)
	if !white {
		value = -value
	}
	return value
}

// returns found doubled pawns as a sum of white and black pawns
// where per doubled white pawn +1 and per black pawn -1
func (l *DarkChess) pawn_structure() int {
	doubled_pawns := 0
	isolated_pawns := 0
	white_pawns_cols := make(map[int]bool)
	black_pawns_cols := make(map[int]bool)
	for col := 0; col < int(row_length_dc); col++ {
		white_pawn_in_col := false
		white_double_in_col := false
		black_pawn_in_col := false
		black_double_in_col := false
		for field := col; field < int(no_fields_dc); field += int(row_length_dc) {
			if field%int(row_length_dc) == 0 {
				white_pawn_in_col = false
				white_double_in_col = false
				black_pawn_in_col = false
				black_double_in_col = false
			}
			if l.white_pawns&(0b1<<field) > 0 {
				white_pawns_cols[col] = true
				if white_pawn_in_col && !white_double_in_col {
					doubled_pawns++
					white_double_in_col = true
				}
				white_pawn_in_col = true
			}
			if l.black_pawns&(0b1<<field) > 0 {
				black_pawns_cols[col] = true
				if black_pawn_in_col && !black_double_in_col {
					doubled_pawns--
					black_double_in_col = true
				}
				black_pawn_in_col = true
			}
		}
	}
	for key, _ := range white_pawns_cols {
		_, left := white_pawns_cols[key-1]
		_, right := white_pawns_cols[key+1]
		if !left && !right {
			isolated_pawns++
		}
	}
	for key, _ := range black_pawns_cols {
		_, left := black_pawns_cols[key-1]
		_, right := black_pawns_cols[key+1]
		if !left && !right {
			isolated_pawns--
		}
	}
	return doubled_pawns + isolated_pawns
}

func (l *DarkChess) String() string {
	field_string := ""
	row_string := ""
	for i := int(row_length_dc*row_length_dc) - 1; i >= 0; i-- {
		field_mask := uint64(0b1 << i)
		if l.white_rooks&field_mask > 0 {
			row_string = "R " + row_string
		} else if l.white_knights&field_mask > 0 {
			row_string = "N " + row_string
		} else if l.white_queen&field_mask > 0 {
			row_string = "Q " + row_string
		} else if l.white_king&field_mask > 0 {
			row_string = "K " + row_string
		} else if l.white_pawns&field_mask > 0 {
			row_string = "P " + row_string
		} else if l.white_bishop&field_mask > 0 {
			row_string = "B " + row_string
		} else if l.black_rooks&field_mask > 0 {
			row_string = "r " + row_string
		} else if l.black_knights&field_mask > 0 {
			row_string = "n " + row_string
		} else if l.black_queen&field_mask > 0 {
			row_string = "q " + row_string
		} else if l.black_king&field_mask > 0 {
			row_string = "k " + row_string
		} else if l.black_pawns&field_mask > 0 {
			row_string = "p " + row_string
		} else if l.black_bishop&field_mask > 0 {
			row_string = "b " + row_string
		} else {
			row_string = "0 " + row_string
		}
		if i%int(row_length_dc) == 0 && i != 0 {
			field_string += row_string + "\n"
			row_string = ""
		}
		if i == 0 {
			field_string += row_string
		}
	}
	return field_string
}

func (l *DarkChess) FENString() string {
	field_string := ""
	row_string := ""
	empty_fields := 0
	for i := int(row_length_dc*row_length_dc) - 1; i >= 0; i-- {
		field_mask := uint64(0b1 << i)
		if (l.white_occupancy|l.black_occupancy)&field_mask > 0 && empty_fields > 0 {
			row_string = strconv.Itoa(empty_fields) + row_string
			empty_fields = 0
		}
		if l.white_rooks&field_mask > 0 {
			row_string = "R" + row_string
		} else if l.white_knights&field_mask > 0 {
			row_string = "N" + row_string
		} else if l.white_queen&field_mask > 0 {
			row_string = "Q" + row_string
		} else if l.white_king&field_mask > 0 {
			row_string = "K" + row_string
		} else if l.white_pawns&field_mask > 0 {
			row_string = "P" + row_string
		} else if l.white_bishop&field_mask > 0 {
			row_string = "B" + row_string
		} else if l.black_rooks&field_mask > 0 {
			row_string = "r" + row_string
		} else if l.black_knights&field_mask > 0 {
			row_string = "n" + row_string
		} else if l.black_queen&field_mask > 0 {
			row_string = "q" + row_string
		} else if l.black_king&field_mask > 0 {
			row_string = "k" + row_string
		} else if l.black_pawns&field_mask > 0 {
			row_string = "p" + row_string
		} else if l.black_bishop&field_mask > 0 {
			row_string = "b" + row_string
		} else {
			empty_fields++
		}
		if i%int(row_length_dc) == 0 {
			if empty_fields > 0 {
				row_string = strconv.Itoa(empty_fields) + row_string
				empty_fields = 0
			}
			if i != 0 {
				field_string += row_string + "/"
				row_string = ""
			}
		}
		if i == 0 {
			field_string += row_string
		}
	}
	field_string += strconv.Itoa(empty_fields)
	return field_string
}
