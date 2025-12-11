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
	black_rooks_moves   [no_fields_lac]uint
	white_knights       uint
	black_knights       uint
	knights_moves       [no_fields_lac]uint
	white_king          uint
	black_king          uint
	king_moves          [no_fields_lac]uint
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

	for i := uint(0); i < no_fields_lac; i++ {
		l.init_pawns(i)
		l.init_rooks(i)
		l.init_knights(i)
		l.init_kings(i)
	}
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

func (l *LosAlamosChess) init_rooks(i uint) {
	col := i % row_length_lac
	row := i / row_length_lac
	// 2x up left
	if col != 0 && row < 4 {
		l.knights_moves[i] += i + 2*row_length_lac - 1
	}
	// 2x up right
	if col != 5 && row < 4 {
		l.knights_moves[i] += i + 2*row_length_lac + 1
	}
	// 2x right up
	if col < 4 && row < 5 {
		l.knights_moves[i] += i + row_length_lac + 2
	}
	// 2x right down
	if col < 4 && row > 0 {
		l.knights_moves[i] += i - row_length_lac + 2
	}
	// 2x down right
	if col < 5 && row > 1 {
		l.knights_moves[i] += i - row_length_lac*2 + 1
	}
	// 2x down left
	if col > 0 && row > 1 {
		l.knights_moves[i] += i - row_length_lac*2 - 1
	}
	// 2x left down
	if col > 1 && row > 0 {
		l.knights_moves[i] += i - row_length_lac - 2
	}
	// 2x left up
	if col > 1 && row < 5 {
		l.knights_moves[i] += i + row_length_lac - 2
	}
}

func (l *LosAlamosChess) init_knights(i uint) {
	// TODO: implement
}

func (l *LosAlamosChess) init_kings(i uint) {
	col := i % row_length_lac
	row := i / row_length_lac
	// up
	if row < 5 {
		l.king_moves[i] += i + row_length_lac
	}
	// up right
	if row < 5 || col < 5 {
		l.king_moves[i] += i + row_length_lac + 1
	}
	// right
	if col < 5 {
		l.king_moves[i] += i + 1
	}
	// right down
	if row > 0 || col < 5 {
		l.king_moves[i] += i - row_length_lac + 1
	}
	// down
	if row > 0 {
		l.king_moves[i] += i - row_length_lac
	}
	// down left
	if row > 0 || col > 0 {
		l.king_moves[i] += i - row_length_lac - 1
	}
	// left
	if col > 0 {
		l.king_moves[i] += i - 1
	}
	// left up
	if row < 5 && col > 0 {
		l.king_moves[i] += i + row_length_lac - 1
	}
}

func (l *LosAlamosChess) ReturnBoard() LosAlamosChess {
	// TODO: implement!
	return *l
}

func (l *LosAlamosChess) GetPreviousBoard() LosAlamosChess {
	// TODO: implement!
	return *l
}

func (l *LosAlamosChess) GetNumberOfMoves() int {
	// TODO: implement!
	return 0
}

func (l *LosAlamosChess) PossibleMoves() []Move {
	// TODO: implement!
	return []Move{}
}

func (l *LosAlamosChess) ExecuteMove(move Move) LosAlamosChess {
	// TODO: implement!
	return *l
}

func (l *LosAlamosChess) CreateView() LosAlamosChess {
	// TODO: implement!
	return l.ReturnBoard()
}

func (l *LosAlamosChess) GameOver() (bool, int) {
	// TODO: implement!
	return false, 0
}

func (l *LosAlamosChess) Heuristic() float64 {
	// TODO: implement!
	return 0.0
}

func (l *LosAlamosChess) String() string {
	// TODO: implement!
	return "Not Implemented!"
}
