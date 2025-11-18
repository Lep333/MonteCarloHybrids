package chess_variation

type Move struct {
	From    int8
	To      int8
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
	Heuristic() float64
	String() string
}
