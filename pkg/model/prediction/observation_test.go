package prediction

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetWordsKnownRatio(t *testing.T) {
	o := Observation{
		sampledWords:          []int{2, 5, 6, 8},
		sampledWordsKnowledge: []bool{true, true, false, false},
	}

	var ratio float32
	var totalCount int

	ratio, totalCount = o.getWordsKnownRatio(0, 10)
	require.Equal(t, true, math.Abs(float64(ratio-0.5)) < 0.001, fmt.Sprintf("Actual %3.f, expected: %3.f", ratio, 0.5))
	require.Equal(t, 4, totalCount)

	ratio, totalCount = o.getWordsKnownRatio(6, 10)
	require.Equal(t, true, math.Abs(float64(ratio-0)) < 0.001, fmt.Sprintf("Actual %3.f, expected: %d", ratio, 0))
	require.Equal(t, 2, totalCount)

	ratio, totalCount = o.getWordsKnownRatio(1, 5)
	require.Equal(t, true, math.Abs(float64(ratio-1)) < 0.001, fmt.Sprintf("Actual %3.f, expected: %d", ratio, 1))
	require.Equal(t, 2, totalCount)

	ratio, totalCount = o.getWordsKnownRatio(12, 15)
	require.Equal(t, true, math.Abs(float64(ratio-0)) < 0.001, fmt.Sprintf("Actual %3.f, expected: %d", ratio, 0))
	require.Equal(t, 0, totalCount)
}

func TestIsWordSampled(t *testing.T) {
	o := Observation{
		sampledWords:          []int{2, 5, 6, 8},
		sampledWordsKnowledge: []bool{true, true, false, false},
	}

	require.Equal(t, true, o.isWordSampled(2))
	require.Equal(t, false, o.isWordSampled(0))
}

func TestTestWords(t *testing.T) {
	o := Observation{
		sampledWords:          []int{2, 5, 6, 8},
		sampledWordsKnowledge: []bool{true, true, false, false},
	}

	verify := func(wordRank int) bool {
		return true
	}

	o.testWords(8, 10, 2, verify)
	require.ElementsMatch(t, []int{2, 5, 6, 8, 9, 10}, o.sampledWords)
	require.ElementsMatch(t, []bool{true, true, false, false, true, true}, o.sampledWordsKnowledge)
}
