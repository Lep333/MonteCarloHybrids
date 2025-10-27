package chess_variation

type Move struct {
	from int8
	to   int8
}

type ChessVariation interface {
	InitGame()
	ReturnBoard() ChessVariation
	PossibleMoves() []Move
	ExecuteMove(move Move) ChessVariation
	GameOver() (bool, int)
	String() string
}
