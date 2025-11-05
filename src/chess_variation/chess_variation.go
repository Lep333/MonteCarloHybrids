package chess_variation

type Move struct {
	from    int8
	to      int8
	Capture bool
}

type ChessVariation interface {
	InitGame()
	ReturnBoard() ChessVariation
	GetPreviousBoard() ChessVariation
	GetNumberOfMoves() int
	PossibleMoves() []Move
	ExecuteMove(move Move) ChessVariation
	CreateView() ChessVariation
	GameOver() (bool, int)
	String() string
}
