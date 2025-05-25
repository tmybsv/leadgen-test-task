package hasher

import (
	"strings"
	"testing"
)

func TestMD5_Hash(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"empty string", "", "d41d8cd98f00b204e9800998ecf8427e"},
		{"simple string", "hello", "5d41402abc4b2a76b9719d911017c592"},
		{"numeric string", "12345678942834", "eeb19b8ea41b8b7a01dc5b1dca6a0216"},
		{"special characters", "#!*^&#^$$$$$$$(*&@!^$&@!#%(*%))", "272b15729bc36f2981e7ac7c5f554df3"},
		{"long string", strings.Repeat("abcd", 1000), "69086d3a75d23ba5b2ce21f7c4c3055b"},
		{"unicode", "привет", "608333adc72f545078ede3aad71bfe74"},
		{"new lines", "adsf\nasfdsadf\nasdfasdfa", "53cf7711769667168378c075784972ec"},
	}

	hasher := &MD5{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasher.Hash(tt.input)

			if got != tt.expect {
				t.Errorf("MD5.Hash(%q) expect %q, got %q", tt.input, tt.expect, got)
			}
		})
	}
}

func TestMD5_Hash_Consistency(t *testing.T) {
	input := "consistency_test"
	hasher := &MD5{}

	hash1 := hasher.Hash(input)
	hash2 := hasher.Hash(input)

	if hash1 != hash2 {
		t.Errorf("MD5 not consistent: %q != %q", hash1, hash2)
	}
}
