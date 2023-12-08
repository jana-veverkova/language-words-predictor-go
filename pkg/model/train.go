package model

import (
	"fmt"
	"math"
	"runtime"
	"sync"

	"github.com/jana-veverkova/language-words-predictor-go/pkg/model/prediction"
	ps "github.com/jana-veverkova/language-words-predictor-go/pkg/populationsample"
	"github.com/pkg/errors"
)

type result struct {
	wordsKnown          int
	lowerGuess          int
	upperGuess          int
	numberOfTestedWords int
}

type summary struct {
	accuracy   float32
	// mean absolute error
	mae                float32
	averageGuessRange  float32
	averageTestedWords float32
}

// iterates through population test set and stores predicted guess range into
func Train() error {
	runtime.GOMAXPROCS(8)

	testSample, err := ps.GetTestSample()
	if err != nil {
		return errors.WithStack(err)
	}

	trainSample, err := ps.GetTrainSample()
	if err != nil {
		return errors.WithStack(err)
	}

	// run for different setting values to compare results
	for i := 100; i <= 1000; i = i + 100 {
		for k := 5; k <= 30; k = k + 5 {

			var wgRun sync.WaitGroup
			var wgSummarize sync.WaitGroup

			chResult := make(chan result)

			setting := prediction.Setting{
				RangeOfWordsToAsk:  i,
				WordTestSize:       k,
				RequiredGuessRange: 300,
				MaxTestedWords:     100,
			}

			wgSummarize.Add(1)
			go summarize(chResult, setting, &wgSummarize)

			for ix, person := range testSample.Table {
				wgRun.Add(1)
				go runRow(person, testSample.RowSums[ix], trainSample, setting, chResult, &wgRun)
			}

			wgRun.Wait()
			close(chResult)
			wgSummarize.Wait()

		}
	}

	return nil
}

func runRow(person []bool, wordsKnown int, trainSample *ps.PopulationSample, setting prediction.Setting, chResult chan result, wg *sync.WaitGroup) {
	defer wg.Done()

	verify := func(wordRank int) bool {
		return person[wordRank]
	}

	lowerGuess, upperGuess, numberOfTestedWords, err := prediction.Predict(trainSample, verify, setting)
	if err != nil {
		fmt.Println(errors.WithStack(err))
	} else {
		personResult := result{
			wordsKnown:          wordsKnown,
			lowerGuess:          lowerGuess,
			upperGuess:          upperGuess,
			numberOfTestedWords: numberOfTestedWords,
		}
		chResult <- personResult
	}
}

func summarize(chResult chan result, setting prediction.Setting, wg *sync.WaitGroup) {
	defer wg.Done()

	count := 0
	hits := 0
	absoluteError := float64(0)
	guessRangeSum := 0
	testedWordsCount := 0

	for result := range chResult {
		count++
		if result.wordsKnown >= result.lowerGuess && result.wordsKnown <= result.upperGuess {
			hits++
		} else if result.wordsKnown < result.lowerGuess {
			absoluteError = absoluteError + math.Abs(float64(result.wordsKnown)-float64(result.lowerGuess))
		} else if result.wordsKnown > result.upperGuess {
			absoluteError = absoluteError + math.Abs(float64(result.wordsKnown)-float64(result.upperGuess))
		}
		guessRangeSum = guessRangeSum + result.upperGuess - result.lowerGuess
		testedWordsCount = testedWordsCount + result.numberOfTestedWords
	}

	summary := summary{
		accuracy:           float32(hits) / float32(count),
		mae:                float32(absoluteError) / float32(count),
		averageGuessRange:  float32(guessRangeSum) / float32(count),
		averageTestedWords: float32(testedWordsCount) / float32(count),
	}

	fmt.Printf("Parameters: RangeOfWordsToAsk: %d, WordTestSize: %d, RequiredGuessRange: %d, MaxTestedWords: %d => ", setting.RangeOfWordsToAsk, setting.WordTestSize, setting.RequiredGuessRange, setting.MaxTestedWords)
	fmt.Printf("Accuracy: %.2f, Mae: %.2f, Average guess range: %.2f, Average tested words: %2.f \n", summary.accuracy, summary.mae, summary.averageGuessRange, summary.averageTestedWords)

}
