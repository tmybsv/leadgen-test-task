package hasher

import (
	"crypto/sha256"
	"fmt"
)

// SHA256 is a SHA256 hasher.
type SHA256 struct{}

// Hash hashes input string to SHA256 checksum.
func (*SHA256) Hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}
