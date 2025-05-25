package hash

import "errors"

// Hash domain errors.
var (
	ErrEmptyInput           = errors.New("input string cannot be empty")
	ErrEmptyHash            = errors.New("hashed string cannot be empty")
	ErrUnsupportedAlgorithm = errors.New("unsupported algorithm")
)

// Hash represents hash domain entity.
type Hash struct {
	input  string
	hashed string
	alg    Algorithm
}

// New creates new hash instance.
func New(input, hashed string, alg Algorithm) (*Hash, error) {
	if input == "" {
		return nil, ErrEmptyInput
	}

	if hashed == "" {
		return nil, ErrEmptyHash
	}

	if !isValidAlgorithm(alg) {
		return nil, ErrUnsupportedAlgorithm
	}

	return &Hash{
		input:  input,
		hashed: hashed,
		alg:    alg,
	}, nil
}

// Hashed returns hashed string.
func (h *Hash) Hashed() string { return h.hashed }

// Algorithm returns hash algorithm
func (h *Hash) Algorithm() Algorithm { return h.alg }

// Input returns string from which hash was build.
func (h *Hash) Input() string { return h.input }

func isValidAlgorithm(alg Algorithm) bool {
	return alg == AlgorithmMD5 || alg == AlgorithmSHA256
}

// Algorithm represents hash algorithm.
type Algorithm int8

// Supported hash algorithms.
const (
	AlgorithmMD5 Algorithm = iota + 1
	AlgorithmSHA256
)

// String strings algorithm numeric constant.
func (a Algorithm) String() string {
	switch a {
	case AlgorithmMD5:
		return "md5"
	case AlgorithmSHA256:
		return "sha256"
	default:
		return ""
	}
}
