package prediction

import (
	"math"

	ps "github.com/jana-veverkova/language-words-predictor-go/pkg/populationsample"
)

const (
	centerRatio = 0.5
	defaultLowerGuess = 0
	defaultUpperGuess = 20000
)

type Setting struct {
	// words will be asked from range 2*RangeOfWordsToAsk
	RangeOfWordsToAsk int
	// in every iteration WordTestSize of words will be asked
	WordTestSize int
	// if guessed range is smalled than RequiredGuessRange iterations stop
	RequiredGuessRange int
	// is number of tested words exceeds MaxTestedWords iterations stop
	MaxTestedWords int
}

type Iteration struct {
	lowerBound int
	upperBound int
	subsample  *Neighbors
	lowerGuess int
	upperGuess int
}

// predicts the number of words an observation knows
// Input:
//// 	populationSample - training set of sampled population to compare observation with
//// 	verify - function the test word and returns weather observation knows or doesn't know the word
// Output: returns lower and higher bound for guess and number of words tested

func Predict(populationSample *ps.PopulationSample, verify func(wordRank int) bool, setting Setting) (int, int, int, error) {
	observation := Observation{
		sampledWords:          make([]int, 0),
		sampledWordsKnowledge: make([]bool, 0),
	}

	// get row indices of population sample
	rowIxs := make([]int, 0)
	for ix := range populationSample.Table {
		rowIxs = append(rowIxs, ix)
	}

	population := Neighbors{populationSample, rowIxs}
	neighbors := Neighbors{populationSample, rowIxs}

	iterations := make([]*Iteration, 0)

	lowerGuess := defaultLowerGuess
	upperGuess := defaultUpperGuess

	counter := 0
	for upperGuess-lowerGuess > setting.RequiredGuessRange && len(observation.sampledWords) < setting.MaxTestedWords {
		counter++
		iteration := Iteration{}

		// find column index centerIx where words known are closest to 0.5
		// get words from centerIx - range and centerIx + range
		_, lowerBound, upperBound := getWordTestBounds(&neighbors, centerRatio, setting)

		// test words from range (lowerBound, upperBound)
		observation.testWords(lowerBound, upperBound, setting.WordTestSize, verify)

		// select person from population training set that know approximately same number of words in the range
		subsample, subsampleLowerGuess, subsampleUpperGuess := iterate(&population, &observation, lowerBound, upperBound)

		iteration.lowerBound = lowerBound
		iteration.upperBound = upperBound
		iteration.subsample = subsample
		iteration.lowerGuess = subsampleLowerGuess
		iteration.upperGuess = subsampleUpperGuess
		iterations = append(iterations, &iteration)

		// reiterate all iterations - consider new tested words
		for _, iter := range iterations {
			(*iter).subsample, (*iter).lowerGuess, (*iter).upperGuess = iterate(&population, &observation, iter.lowerBound, iter.upperBound)
		}

		var finalGroupIxs []int
		finalGroupIxs, lowerGuess, upperGuess = createFinalGroup(iterations)
		neighbors = Neighbors{populationSample, finalGroupIxs}

		// fmt.Printf("Iteration %d \n", counter)
		// fmt.Printf("Word bounds: %d (%d, %d)\n", centerIx, lowerBound, upperBound)
		// fmt.Printf("Words asked: %d\n", len(observation.sampledWords))
		// fmt.Printf("Guess: (%d, %d)\n", lowerGuess, upperGuess)
	}

	return lowerGuess, upperGuess, len(observation.sampledWords), nil
}

// return population subsample and lower and upper guess
func iterate(population *Neighbors, observation *Observation, lowerBound int, upperBound int) (*Neighbors, int, int) {
	// compute confidence intervals
	_, confidenceLower, confidenceUpper := computeConfidenceIntervals(observation, lowerBound, upperBound)

	// select persons having words known ratio inside confidence interval
	subsample := population.filterByConfidenceInterval(lowerBound, upperBound, confidenceLower, confidenceUpper)

	// compute words known by this subsample
	subsampleLowerGuess, subsampleUpperGuess := subsample.getWordsKnownMinMax()

	return subsample, subsampleLowerGuess, subsampleUpperGuess
}

func getWordTestBounds(neighbors *Neighbors, ratio float32, setting Setting) (int, int, int) {
	centerIx := neighbors.getWordRankByKnowledgeAverage(ratio)
	lowerBound := int(math.Max(float64(centerIx-setting.RangeOfWordsToAsk), 0))
	upperBound := int(math.Min(float64(centerIx+setting.RangeOfWordsToAsk), float64(len(neighbors.PopulationSample.Table[0]))))
	return centerIx, lowerBound, upperBound
}

func computeConfidenceIntervals(observation *Observation, lowerBound int, upperBound int) (float32, float32, float32) {
	wordsKnownRatio, sampleSize := observation.getWordsKnownRatio(lowerBound, upperBound)

	sd := math.Sqrt(math.Max(float64(wordsKnownRatio),0.05)*math.Max(float64(1-wordsKnownRatio),0.05)/float64(sampleSize))
	intervalLower := wordsKnownRatio - float32(2*sd)
	intervalUpper := wordsKnownRatio + float32(2*sd)

	return wordsKnownRatio, intervalLower, intervalUpper
}

// final group is intersection of all groups found in all iterations
// output: final group ixs, lowerGuess, upperGues
func createFinalGroup(iterations []*Iteration) ([]int, int, int) {
	var finalGroupIxs []int
	var lowerGuess int
	var upperGuess int

	for ix, iter := range iterations {
		if ix == 0 {
			finalGroupIxs = iter.subsample.Ixs

			lowerGuess = iter.lowerGuess
			upperGuess = iter.upperGuess
			continue
		}
		finalGroupIxs = makeIntersection(finalGroupIxs, iter.subsample.Ixs)
		if iter.lowerGuess > lowerGuess {
			lowerGuess = iter.lowerGuess
		}
		if iter.upperGuess < upperGuess {
			upperGuess = iter.upperGuess
		}
	}

	return finalGroupIxs, lowerGuess, upperGuess
}

func makeIntersection(sample1 []int, sample2 []int) []int {
	intersection := make([]int, 0)

	var found bool

	for _, s1 := range sample1 {
		found = false
		for _, s2 := range sample2 {
			if s2 == s1 {
				found = true
			}
		}
		if found {
			intersection = append(intersection, s1)
		}
	}

	return intersection
}
