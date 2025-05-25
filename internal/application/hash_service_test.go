package application

import (
	"context"
	"errors"
	"testing"

	"github.com/tmybsv/leadgen-test-task/internal/domain/hash"
)

type mockRepository struct {
	findByInputFunc func(ctx context.Context, input string, alg hash.Algorithm) (*hash.Hash, error)
	saveFunc        func(ctx context.Context, h *hash.Hash) error
}

func (m *mockRepository) FindByInput(ctx context.Context, input string, alg hash.Algorithm) (*hash.Hash, error) {
	return m.findByInputFunc(ctx, input, alg)
}

func (m *mockRepository) Save(ctx context.Context, h *hash.Hash) error {
	return m.saveFunc(ctx, h)
}

type mockHasher struct {
	hashFunc func(input string) string
}

func (m *mockHasher) Hash(input string) string {
	return m.hashFunc(input)
}

func TestNewHashService(t *testing.T) {
	repo := &mockRepository{}
	hashers := map[hash.Algorithm]hash.Hasher{
		hash.AlgorithmMD5: &mockHasher{},
	}

	service := NewHashService(repo, hashers)

	if service.hashRepo != repo {
		t.Error("repo not set")
	}

	if len(service.hashers) != 1 {
		t.Error("hashers not set")
	}
}

func TestHashService_CreateHash(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		alg            hash.Algorithm
		repoFindResult *hash.Hash
		repoFindError  error
		hasherExists   bool
		hasherResult   string
		repoSaveError  error
		expectError    bool
	}{
		{
			name:           "cache hit",
			input:          "test",
			alg:            hash.AlgorithmMD5,
			repoFindResult: mustCreateHash("test", "cached_hash", hash.AlgorithmMD5),
			repoFindError:  nil,
			expectError:    false,
		},
		{
			name:          "cache miss success",
			input:         "test",
			alg:           hash.AlgorithmMD5,
			repoFindError: errors.New("not found"),
			hasherExists:  true,
			hasherResult:  "new_hash",
			repoSaveError: nil,
			expectError:   false,
		},
		{
			name:          "hasher not registered",
			input:         "test",
			alg:           hash.AlgorithmSHA256,
			repoFindError: errors.New("not found"),
			hasherExists:  false,
			expectError:   true,
		},
		{
			name:          "save error",
			input:         "test",
			alg:           hash.AlgorithmMD5,
			repoFindError: errors.New("not found"),
			hasherExists:  true,
			hasherResult:  "new_hash",
			repoSaveError: errors.New("save failed"),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{
				findByInputFunc: func(_ context.Context, _ string, _ hash.Algorithm) (*hash.Hash, error) {
					return tt.repoFindResult, tt.repoFindError
				},
				saveFunc: func(_ context.Context, _ *hash.Hash) error {
					return tt.repoSaveError
				},
			}

			hashers := map[hash.Algorithm]hash.Hasher{}
			if tt.hasherExists {
				hashers[tt.alg] = &mockHasher{
					hashFunc: func(_ string) string {
						return tt.hasherResult
					},
				}
			}

			service := NewHashService(repo, hashers)
			result, err := service.CreateHash(context.Background(), tt.input, tt.alg)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
				return
			}

			if !tt.expectError && result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}

func mustCreateHash(input, hashed string, alg hash.Algorithm) *hash.Hash {
	h, err := hash.New(input, hashed, alg)
	if err != nil {
		panic(err)
	}
	return h
}
