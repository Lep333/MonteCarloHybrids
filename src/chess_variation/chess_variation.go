package chess_variation

type Piece int

const (
	Pawn Piece = iota
	Rook
	Knight
	Queen
	King
	Bishop
)

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
	CreateView(white bool) ChessVariation
	GetView(white bool) uint64
	ViewHash(white bool) uint64
	Hash() uint64
	Create_fallback_particle(belief ChessVariation, white bool) ChessVariation
	GameOver() (bool, int)
	Heuristic(white bool) float64
	String() string
	FENString() string
}
