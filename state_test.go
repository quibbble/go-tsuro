package go_tsuro

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewState(t *testing.T) {
	testCases := []struct {
		name      string
		teams     []string
		random    *rand.Rand
		variant   string
		shouldErr bool
	}{
		{
			name:      "teams greater than 11 should error",
			teams:     []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"},
			random:    rand.New(rand.NewSource(123)),
			variant:   VariantClassic,
			shouldErr: true,
		},
		{
			name:      "duplicate teams should error",
			teams:     []string{"1", "2", "1"},
			random:    rand.New(rand.NewSource(123)),
			variant:   VariantClassic,
			shouldErr: true,
		},
		{
			name:      "nil random seed should error",
			teams:     []string{"1", "2"},
			random:    nil,
			variant:   VariantClassic,
			shouldErr: true,
		},
		{
			name:      "invalid variant should error",
			teams:     []string{"1", "2"},
			random:    rand.New(rand.NewSource(123)),
			variant:   "VariantInvalid",
			shouldErr: true,
		},
		{
			name:      "missing variant should error",
			teams:     []string{"1", "2"},
			random:    rand.New(rand.NewSource(123)),
			variant:   "",
			shouldErr: true,
		},
	}
	for _, test := range testCases {
		_, err := newState(test.teams, test.random, test.variant)
		assert.Equal(t, err != nil, test.shouldErr, "ERROR: ", test.name)
	}
}
