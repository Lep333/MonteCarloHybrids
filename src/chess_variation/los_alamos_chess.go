package chess_variation

import (
	"math"
	"math/rand/v2"
	"strconv"
)

const row_length_lac uint = 6
const no_fields_lac uint = row_length_lac * row_length_lac
const black_base_line_start_lac uint = row_length_lac * (row_length_lac - 1)
const no_of_piece_types = 5
const len_zobrist_numbers = no_fields_lac * no_of_piece_types * 2

var lac_white_pawns_moves [no_fields_lac]uint
var lac_white_pawns_capture [no_fields_lac]uint
var lac_black_pawns_moves [no_fields_lac]uint
var lac_black_pawns_capture [no_fields_lac]uint
var lac_white_rooks_moves [no_fields_lac]uint
var lac_knights_moves [no_fields_lac]uint
var lac_king_moves [no_fields_lac]uint
var zobrist_numbers [len_zobrist_numbers]uint64

type LosAlamosChess struct {
	white_pawns     uint
	black_pawns     uint
	white_rooks     uint
	black_rooks     uint
	white_knights   uint
	black_knights   uint
	white_king      uint
	black_king      uint
	white_queen     uint
	black_queen     uint
	white_occupancy uint
	black_occupancy uint
	whiteToPlay     bool
	number_of_moves int
	move_count      int
}

func init() {
	for i := 0; i < int(len_zobrist_numbers); i++ {
		zobrist_numbers[i] = rand.Uint64()
	}
	for i := uint(0); i < no_fields_lac; i++ {
		init_pawns(i)
		init_knights(i)
		init_kings(i)
	}
}

func (l *LosAlamosChess) InitGame() {
	l.whiteToPlay = true
	l.number_of_moves = 0
	l.white_pawns = 0b111111000000
	l.black_pawns = 0b111111000000 << (row_length_lac * 3)
	l.white_rooks = 0b100001
	l.black_rooks = l.white_rooks << (row_length_lac * 5)
	l.white_knights = 0b010010
	l.black_knights = l.white_knights << (row_length_lac * 5)
	l.white_king = 0b001000
	l.black_king = l.white_king << (row_length_lac * 5)
	l.white_queen = 0b00100
	l.black_queen = l.white_queen << (row_length_lac * 5)

	l.set_occupancy_boards()
}

func init_pawns(i uint) {
	if i < no_fields_lac-row_length_lac {
		lac_white_pawns_moves[i] = 0b1 << (i + row_length_lac)
		if i%row_length_lac != row_length_lac-1 {
			lac_white_pawns_capture[i] += 0b1 << (i + row_length_lac + 1) // left capture
		}
		if i%row_length_lac != 0 {
			lac_white_pawns_capture[i] += 0b1 << (i + row_length_lac - 1) // right capture
		}
	}

	if i >= row_length_lac {
		lac_black_pawns_moves[i] = 0b1 << (i - row_length_lac)

		if i%row_length_lac != row_length_lac-1 {
			lac_black_pawns_capture[i] += 0b1 << (i - row_length_lac + 1) // left capture
		}
		if i%row_length_lac != 0 {
			lac_black_pawns_capture[i] += 0b1 << (i - row_length_lac - 1) // right capture
		}
	}
}

func init_knights(i uint) {
	col := i % row_length_lac
	row := i / row_length_lac
	// 2x up left
	if col != 0 && row < 4 {
		lac_knights_moves[i] += 0b1 << (i + 2*row_length_lac - 1)
	}
	// 2x up right
	if col != 5 && row < 4 {
		lac_knights_moves[i] += 0b1 << (i + 2*row_length_lac + 1)
	}
	// 2x right up
	if col < 4 && row < 5 {
		lac_knights_moves[i] += 0b1 << (i + row_length_lac + 2)
	}
	// 2x right down
	if col < 4 && row > 0 {
		lac_knights_moves[i] += 0b1 << (i - row_length_lac + 2)
	}
	// 2x down right
	if col < 5 && row > 1 {
		lac_knights_moves[i] += 0b1 << (i - row_length_lac*2 + 1)
	}
	// 2x down left
	if col > 0 && row > 1 {
		lac_knights_moves[i] += 0b1 << (i - row_length_lac*2 - 1)
	}
	// 2x left down
	if col > 1 && row > 0 {
		lac_knights_moves[i] += 0b1 << (i - row_length_lac - 2)
	}
	// 2x left up
	if col > 1 && row < 5 {
		lac_knights_moves[i] += 0b1 << (i + row_length_lac - 2)
	}
}

func init_kings(i uint) {
	col := i % row_length_lac
	row := i / row_length_lac

	// up
	if row < 5 {
		lac_king_moves[i] += 0b1 << (i + row_length_lac)
	}
	// up right
	if row < 5 && col < 5 {
		lac_king_moves[i] += 0b1 << (i + row_length_lac + 1)
	}
	// right
	if col < 5 {
		lac_king_moves[i] += 0b1 << (i + 1)
	}
	// right down
	if row > 0 && col < 5 {
		lac_king_moves[i] += 0b1 << (i - row_length_lac + 1)
	}
	// down
	if row > 0 {
		lac_king_moves[i] += 0b1 << (i - row_length_lac)
	}
	// down left
	if row > 0 && col > 0 {
		lac_king_moves[i] += 0b1 << (i - row_length_lac - 1)
	}
	// left
	if col > 0 {
		lac_king_moves[i] += 0b1 << (i - 1)
	}
	// left up
	if row < 5 && col > 0 {
		lac_king_moves[i] += 0b1 << (i + row_length_lac - 1)
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
	moves := [200]Move{}
	l.set_occupancy_boards()
	l.move_count = l.generate_moves(moves[:])
	return moves[:l.move_count]
}

func (l *LosAlamosChess) generate_moves(buffer []Move) int {
	n := 0
	own_pawns := l.white_pawns
	own_pawns_moves := lac_white_pawns_moves
	own_pawns_capture := lac_white_pawns_capture
	own_rooks := l.white_rooks
	own_knights := l.white_knights
	own_queen := l.white_queen
	own_king := l.white_king
	own_occupancy := l.white_occupancy
	opponent_occupancy := l.black_occupancy
	white_to_play := true
	if !l.whiteToPlay {
		own_pawns = l.black_pawns
		own_pawns_moves = lac_black_pawns_moves
		own_pawns_capture = lac_black_pawns_capture
		own_rooks = l.black_rooks
		own_knights = l.black_knights
		own_queen = l.black_queen
		own_king = l.black_king
		own_occupancy = l.black_occupancy
		opponent_occupancy = l.white_occupancy
		white_to_play = false
	}

	for i := uint(0); i < no_fields_lac; i++ {
		// pawns
		if own_pawns&(0b1<<i) > 0 {
			moves_possible := (own_pawns_moves[i] & ^opponent_occupancy &
				^own_occupancy) |
				(own_pawns_capture[i] & opponent_occupancy)
			n = l.move_bitboard_to_moves(i, moves_possible, buffer, n, opponent_occupancy)
		}
		// rooks
		if own_rooks&(0b1<<i) > 0 {
			n = l.get_rook_moves(i, white_to_play, buffer, n)
		}
		// knights
		if own_knights&(0b1<<i) > 0 {
			moves_possible := lac_knights_moves[i] & ^own_occupancy
			n = l.move_bitboard_to_moves(i, moves_possible, buffer, n, opponent_occupancy)
		}
		// queen
		if own_queen&(0b1<<i) > 0 {
			n = l.get_rook_moves(i, white_to_play, buffer, n)
			n = l.get_bishop_moves(i, white_to_play, buffer, n)
		}
		// king
		if own_king&(0b1<<i) > 0 {
			moves_possible := lac_king_moves[i] & (^own_occupancy)
			n = l.move_bitboard_to_moves(i, moves_possible, buffer, n, opponent_occupancy)
		}
	}
	return n
}

func (l *LosAlamosChess) move_bitboard_to_moves(start uint, move_bitboard uint, buffer []Move, n int, opp_occupancy uint) int {
	var piece_type Piece
	for i := uint(0); i < no_fields_lac; i++ {
		move_to := uint(0b1 << i)
		if move_bitboard&move_to > 0 {
			capture := false
			if opp_occupancy&move_to > 0 {
				capture = true
				piece_type = l.get_captured_piece(move_to)
			}
			move := Move{From: int8(start), To: int8(i), Capture: capture, CapturedPiece: piece_type}
			buffer[n] = move
			n++
		}
	}
	return n
}

func (l *LosAlamosChess) get_captured_piece(move_to uint) Piece {
	var piece_type Piece
	opp_king := l.black_king
	opp_queen := l.black_queen
	opp_pawn := l.black_pawns
	opp_rook := l.black_rooks
	opp_knights := l.black_knights
	if !l.whiteToPlay {
		opp_king = l.white_king
		opp_queen = l.white_queen
		opp_pawn = l.white_pawns
		opp_rook = l.white_rooks
		opp_knights = l.white_knights
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
	}
	return piece_type
}

func (l *LosAlamosChess) get_rook_moves(i uint, white_to_play bool, buffer []Move, n int) int {
	var piece_type Piece
	own_occupancy := l.white_occupancy
	opponent_occupancy := l.black_occupancy
	if !white_to_play {
		own_occupancy = l.black_occupancy
		opponent_occupancy = l.white_occupancy
	}
	// up
	index := int(i + row_length_lac)
	for index < int(no_fields_lac) && !(i/row_length_lac == row_length_lac-1) {
		move_to := uint(0b1 << index)
		if own_occupancy&(move_to) > 0 {
			break
		}
		if opponent_occupancy&(move_to) > 0 {
			piece_type = l.get_captured_piece(move_to)
			move := Move{From: int8(i), To: int8(index), Capture: true, CapturedPiece: piece_type}
			buffer[n] = move
			n++
			break
		}
		move := Move{From: int8(i), To: int8(index), Capture: false}
		buffer[n] = move
		n++
		index += int(row_length_lac)
	}
	// down
	index = int(i) - int(row_length_lac)
	for index >= 0 && !(i/row_length_lac == 0) {
		move_to := uint(0b1 << index)
		if own_occupancy&(move_to) > 0 {
			break
		}
		if opponent_occupancy&(0b1<<index) > 0 {
			piece_type = l.get_captured_piece(move_to)
			move := Move{From: int8(i), To: int8(index), Capture: true, CapturedPiece: piece_type}
			buffer[n] = move
			n++
			break
		}
		move := Move{From: int8(i), To: int8(index), Capture: false}
		buffer[n] = move
		n++
		index -= int(row_length_lac)
	}
	// right
	index = int(i) + 1
	for index < int(no_fields_lac) && !(i%row_length_lac == 5) {
		capture := false
		move_to := uint(0b1 << index)
		if own_occupancy&(move_to) > 0 {
			break
		}
		if opponent_occupancy&(move_to) > 0 {
			capture = true
			piece_type = l.get_captured_piece(move_to)
		}
		if index%int(row_length_lac) == int(row_length_lac)-1 || capture {
			move := Move{From: int8(i), To: int8(index), Capture: capture, CapturedPiece: piece_type}
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
	for index >= 0 && !(i%row_length_lac == 0) {
		capture := false
		move_to := uint(0b1 << index)
		if own_occupancy&(move_to) > 0 {
			break
		}
		if opponent_occupancy&(move_to) > 0 {
			capture = true
			piece_type = l.get_captured_piece(move_to)
		}
		if index%int(row_length_lac) == 0 || capture {
			move := Move{From: int8(i), To: int8(index), Capture: capture, CapturedPiece: piece_type}
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

func (l *LosAlamosChess) get_bishop_moves(i uint, white_to_play bool, buffer []Move, n int) int {
	var piece_type Piece
	own_occupancy := l.white_occupancy
	opponent_occupancy := l.black_occupancy
	if !white_to_play {
		own_occupancy = l.black_occupancy
		opponent_occupancy = l.white_occupancy
	}
	// up right
	index := int(i + 7)
	for index < int(no_fields_lac) && !(i%row_length_lac == 5) {
		capture := false
		move_to := uint(0b1 << index)
		if own_occupancy&(move_to) > 0 {
			break
		}
		if opponent_occupancy&move_to > 0 {
			capture = true
			piece_type = l.get_captured_piece(move_to)
		}
		move := Move{From: int8(i), To: int8(index), Capture: capture, CapturedPiece: piece_type}
		buffer[n] = move
		n++
		if index%int(row_length_lac) == int(row_length_lac)-1 || capture {
			break
		}
		index += 7
	}
	// down right
	index = int(i) - 5
	for index >= 0 && !(i%row_length_lac == 5) {
		capture := false
		move_to := uint(0b1 << index)
		if own_occupancy&move_to > 0 {
			break
		}
		if opponent_occupancy&move_to > 0 {
			capture = true
			piece_type = l.get_captured_piece(move_to)
		}
		move := Move{From: int8(i), To: int8(index), Capture: capture, CapturedPiece: piece_type}
		buffer[n] = move
		n++
		if index%int(row_length_lac) == int(row_length_lac)-1 || capture {
			break
		}
		index -= 5
	}
	// down left
	index = int(i) - 7
	for index >= 0 && !(i%row_length_lac == 0) {
		capture := false
		move_to := uint(0b1 << index)
		if own_occupancy&move_to > 0 {
			break
		}
		if opponent_occupancy&move_to > 0 {
			capture = true
			piece_type = l.get_captured_piece(move_to)
		}
		move := Move{From: int8(i), To: int8(index), Capture: capture, CapturedPiece: piece_type}
		buffer[n] = move
		n++
		if index%int(row_length_lac) == 0 || capture {
			break
		}
		index -= 7
	}
	// up left
	index = int(i) + 5
	for index < int(no_fields_lac) && !(i%row_length_lac == 0) {
		capture := false
		move_to := uint(0b1 << index)
		if own_occupancy&move_to > 0 {
			break
		}
		if opponent_occupancy&move_to > 0 {
			capture = true
			piece_type = l.get_captured_piece(move_to)
		}
		move := Move{From: int8(i), To: int8(index), Capture: capture, CapturedPiece: piece_type}
		buffer[n] = move
		n++
		if index%int(row_length_lac) == 0 || capture {
			break
		}
		index += 5
	}
	return n
}

func (l *LosAlamosChess) compute_vision(white bool) uint {
	flip_players_turn := false
	local_moves := [200]Move{}
	vision := uint(0)
	if white != l.whiteToPlay {
		l.whiteToPlay = !l.whiteToPlay
		flip_players_turn = true
	}
	move_count := l.generate_moves(local_moves[:])
	for _, move := range local_moves[:move_count] {
		vision |= (1 << move.To) | (1 << move.From)
	}
	if white {
		field_in_front_of_pawns := l.white_pawns << row_length_lac
		vision |= l.white_pawns | l.white_rooks | l.white_knights |
			l.white_queen | l.white_king | field_in_front_of_pawns
	} else {
		field_in_front_of_pawns := l.black_pawns >> row_length_lac
		vision |= l.black_pawns | l.black_rooks | l.black_knights |
			l.black_queen | l.black_king | field_in_front_of_pawns
	}
	if flip_players_turn {
		l.whiteToPlay = !l.whiteToPlay
	}
	return vision
}

func (l *LosAlamosChess) ExecuteMove(move Move) {
	move_to_mask := uint(0b1 << move.To)
	move_from_mask := uint(0b1 << move.From)
	if l.white_rooks&move_to_mask > 0 {
		l.white_rooks = l.white_rooks &^ move_to_mask
	} else if l.white_knights&move_to_mask > 0 {
		l.white_knights = l.white_knights &^ move_to_mask
	} else if l.white_queen&move_to_mask > 0 {
		l.white_queen = l.white_queen &^ move_to_mask
	} else if l.white_king&move_to_mask > 0 {
		l.white_king = l.white_king &^ move_to_mask
	} else if l.white_pawns&move_to_mask > 0 {
		l.white_pawns = l.white_pawns &^ move_to_mask
	} else if l.black_rooks&move_to_mask > 0 {
		l.black_rooks = l.black_rooks &^ move_to_mask
	} else if l.black_knights&move_to_mask > 0 {
		l.black_knights = l.black_knights &^ move_to_mask
	} else if l.black_queen&move_to_mask > 0 {
		l.black_queen = l.black_queen &^ move_to_mask
	} else if l.black_king&move_to_mask > 0 {
		l.black_king = l.black_king &^ move_to_mask
	} else if l.black_pawns&move_to_mask > 0 {
		l.black_pawns = l.black_pawns &^ move_to_mask
	}

	if l.white_rooks&move_from_mask > 0 {
		l.white_rooks = (l.white_rooks &^ move_from_mask) | move_to_mask
	} else if l.white_knights&move_from_mask > 0 {
		l.white_knights = (l.white_knights &^ move_from_mask) | move_to_mask
	} else if l.white_queen&move_from_mask > 0 {
		l.white_queen = (l.white_queen &^ move_from_mask) | move_to_mask
	} else if l.white_king&move_from_mask > 0 {
		l.white_king = (l.white_king &^ move_from_mask) | move_to_mask
	} else if l.white_pawns&move_from_mask > 0 {
		if move.To >= 30 {
			l.white_queen = l.white_queen | move_to_mask
			l.white_pawns = l.white_pawns &^ move_from_mask
		} else {
			l.white_pawns = (l.white_pawns &^ move_from_mask) | move_to_mask
		}
	} else if l.black_rooks&move_from_mask > 0 {
		l.black_rooks = (l.black_rooks &^ move_from_mask) | move_to_mask
	} else if l.black_knights&move_from_mask > 0 {
		l.black_knights = (l.black_knights &^ move_from_mask) | move_to_mask
	} else if l.black_queen&move_from_mask > 0 {
		l.black_queen = (l.black_queen &^ move_from_mask) | move_to_mask
	} else if l.black_king&move_from_mask > 0 {
		l.black_king = (l.black_king &^ move_from_mask) | move_to_mask
	} else if l.black_pawns&move_from_mask > 0 {
		if move.To <= 5 {
			l.black_queen = l.black_queen | move_to_mask
			l.black_pawns = l.black_pawns &^ move_from_mask
		} else {
			l.black_pawns = (l.black_pawns &^ move_from_mask) | move_to_mask
		}
	}
	l.number_of_moves++
	l.whiteToPlay = !l.whiteToPlay

	l.set_occupancy_boards()
}

func (l *LosAlamosChess) UndoMove(move Move) {
	l.number_of_moves--
	l.whiteToPlay = !l.whiteToPlay

	move_to_mask := uint(0b1 << move.To)
	move_from_mask := uint(0b1 << move.From)

	// undo move
	if l.white_pawns&move_to_mask > 0 {
		l.white_pawns = (l.white_pawns &^ move_to_mask) | move_from_mask
		if move.To >= 30 {
			l.white_queen = l.white_queen &^ move_to_mask
		}
	} else if l.white_rooks&move_to_mask > 0 {
		l.white_rooks = (l.white_rooks &^ move_to_mask) | move_from_mask
	} else if l.white_knights&move_to_mask > 0 {
		l.white_knights = (l.white_knights &^ move_to_mask) | move_from_mask
	} else if l.white_queen&move_to_mask > 0 {
		l.white_queen = (l.white_queen &^ move_to_mask) | move_from_mask
	} else if l.white_king&move_from_mask > 0 {
		l.white_king = (l.white_king &^ move_to_mask) | move_from_mask
	} else if l.black_pawns&move_to_mask > 0 {
		l.black_pawns = (l.black_pawns &^ move_to_mask) | move_from_mask
		if move.To <= 5 {
			l.black_queen = l.black_queen &^ move_to_mask
		}
	} else if l.black_rooks&move_to_mask > 0 {
		l.black_rooks = (l.black_rooks &^ move_to_mask) | move_from_mask
	} else if l.black_knights&move_to_mask > 0 {
		l.black_knights = (l.black_knights &^ move_to_mask) | move_from_mask
	} else if l.black_queen&move_to_mask > 0 {
		l.black_queen = (l.black_queen &^ move_to_mask) | move_from_mask
	} else if l.black_king&move_to_mask > 0 {
		l.black_king = (l.black_king &^ move_to_mask) | move_from_mask
	}

	// undo capture
	if l.whiteToPlay {
		switch move.CapturedPiece {
		case Pawn:
			l.black_pawns |= move_to_mask
		case Rook:
			l.black_rooks |= move_to_mask
		case Knight:
			l.black_knights |= move_to_mask
		case Queen:
			l.black_queen |= move_to_mask
		case King:
			l.black_king |= move_to_mask
		}
	} else {
		switch move.CapturedPiece {
		case Pawn:
			l.white_pawns |= move_to_mask
		case Rook:
			l.white_rooks |= move_to_mask
		case Knight:
			l.white_knights |= move_to_mask
		case Queen:
			l.white_queen |= move_to_mask
		case King:
			l.white_king |= move_to_mask
		}
	}
	l.set_occupancy_boards()
}

func (l *LosAlamosChess) set_occupancy_boards() {
	l.white_occupancy = l.white_rooks | l.white_knights |
		l.white_queen | l.white_king | l.white_pawns
	l.black_occupancy = l.black_rooks | l.black_knights |
		l.black_queen | l.black_king | l.black_pawns
}

var moves = [200]Move{}

func (l *LosAlamosChess) CreateView(white bool) ChessVariation {
	copy := *l
	moves := l.PossibleMoves()
	vision := uint(0)
	for _, move := range moves {
		vision |= (1 << move.To)
	}
	if white {
		field_in_front_of_pawns := l.white_pawns << row_length_lac
		vision |= l.white_pawns | l.white_rooks | l.white_knights |
			l.white_queen | l.white_king | field_in_front_of_pawns
	} else {
		field_in_front_of_pawns := l.black_pawns >> row_length_lac
		vision |= l.black_pawns | l.black_rooks | l.black_knights |
			l.black_queen | l.black_king | field_in_front_of_pawns
	}
	copy.white_pawns &= vision
	copy.white_rooks &= vision
	copy.white_knights &= vision
	copy.white_queen &= vision
	copy.white_king &= vision
	copy.black_pawns &= vision
	copy.black_rooks &= vision
	copy.black_knights &= vision
	copy.black_queen &= vision
	copy.black_king &= vision
	copy.set_occupancy_boards()

	return &copy
}

func (l *LosAlamosChess) ViewHash(white bool) uint64 {
	hash := uint64(0)
	mask := l.compute_vision(white)
	gap := no_of_piece_types * 2

	for i := 0; i < int(no_fields_lac); i++ {
		if l.white_pawns&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap]
		} else if l.white_rooks&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+int(Rook)]
		} else if l.white_knights&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+int(Knight)]
		} else if l.white_queen&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+int(Queen)]
		} else if l.white_king&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+int(King)]
		} else if l.black_pawns&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types]
		} else if l.black_rooks&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types+int(Rook)]
		} else if l.black_knights&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types+int(Knight)]
		} else if l.black_queen&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types+int(Queen)]
		} else if l.black_king&mask&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types+int(King)]
		}
	}
	return hash
}

func (l *LosAlamosChess) Hash() uint64 {
	hash := uint64(0)
	gap := no_of_piece_types * 2

	for i := 0; i < int(no_fields_lac); i++ {
		if l.white_pawns&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap]
		} else if l.white_rooks&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+int(Rook)]
		} else if l.white_knights&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+int(Knight)]
		} else if l.white_queen&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+int(Queen)]
		} else if l.white_king&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+int(King)]
		} else if l.black_pawns&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types]
		} else if l.black_rooks&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types+int(Rook)]
		} else if l.black_knights&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types+int(Knight)]
		} else if l.black_queen&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types+int(Queen)]
		} else if l.black_king&(1<<i) > 0 {
			hash ^= zobrist_numbers[i*gap+no_of_piece_types+int(King)]
		}
	}
	return hash
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

	max_moves := 100
	if l.number_of_moves > max_moves {
		return true, 0
	}

	return false, 0
}

func (l *LosAlamosChess) Heuristic(white bool) float64 {
	over, val := l.GameOver()
	if over {
		if white {
			return float64(val)
		} else {
			return float64(-val)
		}
	}
	value := 0.0
	white_material := 0.0
	black_material := 0.0
	for i := 0; i < int(no_fields_lac); i++ {
		if l.white_queen&(0b1<<i) > 0 {
			white_material += 9
		} else if l.white_rooks&(0b1<<i) > 0 {
			white_material += 5
		} else if l.white_knights&(0b1<<i) > 0 {
			white_material += 3
		} else if l.white_pawns&(0b1<<i) > 0 {
			white_material += 1
		} else if l.black_queen&(0b1<<i) > 0 {
			black_material += 9
		} else if l.black_rooks&(0b1<<i) > 0 {
			black_material += 5
		} else if l.black_knights&(0b1<<i) > 0 {
			black_material += 3
		} else if l.black_pawns&(0b1<<i) > 0 {
			black_material += 1
		}
	}
	white_mobility := 0
	black_mobility := 0
	if white {
		white_mobility = len(l.PossibleMoves())
		l.whiteToPlay = !l.whiteToPlay
		black_mobility = len(l.PossibleMoves())
	} else {
		black_mobility = len(l.PossibleMoves())
		l.whiteToPlay = !l.whiteToPlay
		white_mobility = len(l.PossibleMoves())
	}
	l.whiteToPlay = !l.whiteToPlay
	avg_material := (white_material + black_material) / 2
	value = (white_material - black_material + float64(l.pawn_structure()) + 0.1*float64(white_mobility-black_mobility)) / avg_material
	value = math.Min(math.Min(value, -1), 1)
	if !white {
		value = -value
	}
	return value
}

// returns found doubled pawns as a sum of white and black pawns
// where per doubled white pawn +1 and per black pawn -1
func (l *LosAlamosChess) pawn_structure() int {
	doubled_pawns := 0
	isolated_pawns := 0
	white_pawns_cols := make(map[int]bool)
	black_pawns_cols := make(map[int]bool)
	for col := 0; col < int(row_length_lac); col++ {
		for field := col; field < int(no_fields_lac); field += int(row_length_lac) {
			white_pawn_in_col := false
			white_double_in_col := false
			black_pawn_in_col := false
			black_double_in_col := false
			if field%int(row_length_lac) == 0 {
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

func (l *LosAlamosChess) String() string {
	field_string := ""
	row_string := ""
	for i := int(row_length_lac*row_length_lac) - 1; i >= 0; i-- {
		field_mask := uint(0b1 << i)
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
		} else {
			row_string = "0 " + row_string
		}
		if i%int(row_length_lac) == 0 && i != 0 {
			field_string += row_string + "\n"
			row_string = ""
		}
		if i == 0 {
			field_string += row_string
		}
	}
	return field_string
}

func (l *LosAlamosChess) FENString() string {
	field_string := ""
	row_string := ""
	empty_fields := 0
	for i := int(row_length_lac*row_length_lac) - 1; i >= 0; i-- {
		field_mask := uint(0b1 << i)
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
		} else {
			empty_fields++
		}
		if i%int(row_length_lac) == 0 {
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
