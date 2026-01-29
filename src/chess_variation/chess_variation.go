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
	ExecuteMove(move Move)
	UndoMove(move Move)
	CreateView(white bool) ChessVariation
	ViewHash(white bool) uint64
	Hash() uint64
	GameOver() (bool, int)
	Heuristic(white bool) float64
	String() string
	FENString() string
}
