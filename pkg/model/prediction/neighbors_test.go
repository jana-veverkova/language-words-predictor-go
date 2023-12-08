package prediction

import (
	"fmt"
	"math"
	"testing"

	"github.com/jana-veverkova/language-words-predictor-go/pkg/populationsample"
	"github.com/stretchr/testify/require"
)

func TestGetWordKnowledgeAverage(t *testing.T) {
	populationSample := populationsample.PopulationSample{
		Table: [][]bool{
			{true, false, false, false, false, false},
			{true, true, false, false, false, false},
			{true, true, true, false, false, false},
			{true, true, true, true, false, false},
			{true, true, true, true, true, false},
			{true, true, true, true, true, true},
		},
	}

	testData := []Neighbors{
		{PopulationSample: &populationSample, Ixs: []int{0, 1, 2}},
		{PopulationSample: &populationSample, Ixs: []int{1, 2, 3, 5}},
		{PopulationSample: &populationSample, Ixs: []int{0}},
		{PopulationSample: &populationSample, Ixs: []int{}},
	}

	expected := [][]float32{
		{1, 0.666, 0.333, 0, 0, 0},
		{1, 1, 0.75, 0.5, 0.25, 0.25},
		{1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	}

	for i, td := range testData {
		actual := td.getWordKnowledgeAverage()
		for j, act := range actual {
			require.Equal(t, true, math.Abs(float64(expected[i][j])-float64(act)) <= 0.001, fmt.Sprintf("Expected: %4.f, actual: %4.f", float64(expected[i][j]), float64(act)))
		}
	}
}

func TestGetWordKnownMinMax(t *testing.T) {
	populationSample := populationsample.PopulationSample{
		Table: [][]bool{
			{true, false, false, false, false, false},
			{true, true, false, false, false, false},
			{true, true, true, false, false, false},
			{true, true, true, true, false, false},
			{true, true, true, true, true, false},
			{true, true, true, true, true, true},
		},
		RowSums: []int{1, 2, 3, 4, 5, 6},
	}

	testData := []Neighbors{
		{PopulationSample: &populationSample, Ixs: []int{0, 1, 2}},
		{PopulationSample: &populationSample, Ixs: []int{1, 2, 3, 5}},
		{PopulationSample: &populationSample, Ixs: []int{0}},
		{PopulationSample: &populationSample, Ixs: []int{}},
	}

	expected := [][]int{
		{1, 3},
		{2, 6},
		{1, 1},
		{0, 0},
	}

	for i, td := range testData {
		min, max := td.getWordsKnownMinMax()
		require.Equal(t, expected[i][0], min)
		require.Equal(t, expected[i][1], max)
	}
}
