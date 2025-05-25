package hasher

import (
	"crypto/md5"
	"fmt"
)

// MD5 is a MD5 hasher.
type MD5 struct{}

// Hash hashes input string to MD5 checksum.
func (*MD5) Hash(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}
