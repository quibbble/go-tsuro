package go_tsuro

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Duplicates(t *testing.T) {
	dups := []string{"a", "b", "a"}
	assert.True(t, duplicates(dups))

	nondups := []string{"a", "b", "c"}
	assert.False(t, duplicates(nondups))
}
