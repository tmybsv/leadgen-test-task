package hash

// Hasher represents contract that different hash creators should implement.
type Hasher interface {
	Hash(input string) string
}
