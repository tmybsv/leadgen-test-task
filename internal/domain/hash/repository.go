package hash

import "context"

// Repository is a contract that hash repositories should implement.
type Repository interface {
	// Save saves hash.
	Save(context.Context, *Hash) error

	// FindByInput finds hash by input string and algorithm.
	FindByInput(ctx context.Context, input string, alg Algorithm) (*Hash, error)
}
