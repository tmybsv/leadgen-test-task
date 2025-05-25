package hasher

import (
	"strings"
	"testing"
)

func TestSHA256_Hash(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"empty string", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"simple string", "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"numeric string", "12345678942834", "6ccd25f9231059bc0abfcaff987b6327cae3d8f72b0e548172b71738f343d268"},
		{"special characters", "#!*^&#^$$$$$$$(*&@!^$&@!#%(*%))", "6ada9bd72410c28a44d858187a913f351437116753f11c6591562aefdbb75206"},
		{"long string", strings.Repeat("abcd", 1000), "a794d3322e58f529258e0d9331aa9a17eece6e05848aa89fe1e6ffa006231957"},
		{"unicode", "привет", "e58f1e8c55fa105bdd3f40e5037eb0b039b5998d52c05e6cd98878dd2da5cab2"},
		{"new lines", "adsf\nasfdsadf\nasdfasdfa", "3df1e9d8df767f01885455e891152467e3dd682451b17feacb3c575e57eb0b11"},
	}

	hasher := &SHA256{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasher.Hash(tt.input)

			if got != tt.expect {
				t.Errorf("SHA256.Hash(%q) expect %q, got %q", tt.input, tt.expect, got)
			}
		})
	}
}

func TestSHA256_Hash_Consistency(t *testing.T) {
	input := "consistency_test"
	hasher := &SHA256{}

	hash1 := hasher.Hash(input)
	hash2 := hasher.Hash(input)

	if hash1 != hash2 {
		t.Errorf("SHA256 not consistent: %q != %q", hash1, hash2)
	}
}
