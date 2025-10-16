package chess_variation

type Move struct {
	from int8
	to   int8
}

type ChessVariation interface {
	InitGame()
	ReturnBoard() ChessVariation
	PossibleMoves(whiteToPlay bool) []Move
	ExecuteMove(whiteToPlay bool, move Move)
	GameOver() bool
	String() string
}
