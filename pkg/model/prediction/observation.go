package prediction

import (
	"math/rand"
)

// observation is used to store info about tested person
// sampledWords - words already tested
// sampledWordsKnowledge - do the observation know sampled words

type Observation struct {
	sampledWords          []int
	sampledWordsKnowledge []bool
}

func (o *Observation) getWordsKnownRatio(lowerBound int, upperBound int) (float32, int) {
	knownWords := 0
	totalCount := 0
	for ix, b := range o.sampledWordsKnowledge {
		if o.sampledWords[ix] >= lowerBound && o.sampledWords[ix] <= upperBound {
			totalCount++
			if b {
				knownWords++
			}
		}
	}

	if totalCount == 0 {
		return 0, 0
	}
	return float32(knownWords) / float32(totalCount), totalCount
}

func (o *Observation) isWordSampled(wordRank int) bool {
	for _, x := range o.sampledWords {
		if x == wordRank {
			return true
		}
	}
	return false
}

func (o *Observation) testWords(lowerBound int, upperBound int, testSize int, verify func(wordRank int) bool) {
	count := 0

	for i := 0; count != testSize; i++ {
		s := rand.Intn(upperBound+1-lowerBound) + lowerBound
		if o.isWordSampled(s) {
			continue
		}
		count++
		o.sampledWords = append(o.sampledWords, s)
		o.sampledWordsKnowledge = append(o.sampledWordsKnowledge, verify(s))
	}
}
