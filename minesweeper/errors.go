package minesweeper

type GameOverError struct{}

func (*GameOverError) Error() string {
	return "operation failed because game is over"
}

type OutOfFlagsError struct{}

func (*OutOfFlagsError) Error() string {
	return "you are out of flags"
}

type InvalidCoordinatesError struct{}

func (*InvalidCoordinatesError) Error() string {
	return "cell coordinates invalid"
}

type AlreadyFlaggedError struct{}

func (*AlreadyFlaggedError) Error() string {
	return "cell is already flagged"
}

type AlreadyUnflaggedError struct{}

func (*AlreadyUnflaggedError) Error() string {
	return "cell is already unflagged"
}

type AlreadyDuggedError struct{}

func (*AlreadyDuggedError) Error() string {
	return "cell is already dugged"
}
