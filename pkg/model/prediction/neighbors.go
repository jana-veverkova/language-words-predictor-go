package prediction

import (
	"math"

	ps "github.com/jana-veverkova/language-words-predictor-go/pkg/populationsample"
)

// stores subsample of population that is close to observed person

type Neighbors struct {
	PopulationSample *ps.PopulationSample
	Ixs              []int
}

func (n *Neighbors) getWordKnowledgeAverage() []float32 {
	var person []bool

	colSumns := make([]int, len(n.PopulationSample.Table[0]))
	for _, ix := range n.Ixs {
		person = n.PopulationSample.Table[ix]
		for wordRank, isKnown := range person {
			if isKnown {
				colSumns[wordRank]++
			}

		}
	}

	colAverages := make([]float32, len(n.PopulationSample.Table[0]))
	if len(n.Ixs) == 0 {
		return colAverages
	}

	for ix, s := range colSumns {
		colAverages[ix] = float32(s) / float32(len(n.Ixs))
	}

	return colAverages
}

func (n *Neighbors) getWordRankByKnowledgeAverage(average float32) int {
	knowledgeAverage := n.getWordKnowledgeAverage()

	closestRank := 0
	for i, av := range knowledgeAverage {
		if math.Abs(float64(av-average)) < math.Abs(float64(knowledgeAverage[closestRank]-average)) {
			closestRank = i
		}
	}

	return closestRank
}

func (n *Neighbors) getWordsKnownMinMax() (int, int) {
	if len(n.Ixs) == 0 {
		return 0, 0
	}

	min := 100000
	max := 0

	for _, ix := range n.Ixs {
		if n.PopulationSample.RowSums[ix] < min {
			min = n.PopulationSample.RowSums[ix]
		}
		if n.PopulationSample.RowSums[ix] > max {
			max = n.PopulationSample.RowSums[ix]
		}
	}

	return min, max
}

func (n *Neighbors) filterByConfidenceInterval(wordsLowerBound int, wordsUpperBound int, confIntLower float32, confIntUpper float32) *Neighbors {
	filtered := make([]int, 0)

	wordsCount := wordsUpperBound - wordsLowerBound + 1

	var sum int
	for _, ix := range n.Ixs {
		sum = 0
		for c := wordsLowerBound; c <= wordsUpperBound; c++ {
			if n.PopulationSample.Table[ix][c] {
				sum++
			}
		}

		if float32(sum)/float32(wordsCount) >= confIntLower && float32(sum)/float32(wordsCount) <= confIntUpper {
			filtered = append(filtered, ix)
		}
	}

	return &Neighbors{n.PopulationSample, filtered}
}
