package eztools

import "errors"

const (
	// defValidID is the smallest valid ID, different from the default
	defValidID = DefID + 1
	// DefID is the default ID
	DefID = 0
	// AllID stands for all items
	AllID = DefID - 1
	// InvalidID better to be negative to be different from a normal ID
	InvalidID = DefID - 2 //pairs defined some related
)

var (
	// ErrNoValidResults stands for no valid results
	ErrNoValidResults = errors.New("No Valid results")
	// ErrOutOfBound stands for out of bound
	ErrOutOfBound = errors.New("Out of bound")
	// ErrInvalidInput stands for invalid input
	ErrInvalidInput = errors.New("Invalid input")
	// ErrInExistence stands for result for input already in existence
	ErrInExistence = errors.New("In Existence")
)
