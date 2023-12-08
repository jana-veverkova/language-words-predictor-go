package prediction

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeIntersection(t *testing.T) {
	require.ElementsMatch(t, []int{2, 3}, makeIntersection([]int{1, 2, 3}, []int{2, 3, 4}))
	require.ElementsMatch(t, []int{}, makeIntersection([]int{1, 2, 3}, []int{4, 5, 6}))
	require.ElementsMatch(t, []int{}, makeIntersection([]int{}, []int{4, 5, 6}))
	require.ElementsMatch(t, []int{2, 3, 4}, makeIntersection([]int{2, 3, 4}, []int{1, 2, 3, 4, 5, 6}))
}
