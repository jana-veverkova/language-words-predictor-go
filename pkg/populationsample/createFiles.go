package populationsample

import (
	"fmt"
	"math"
	"runtime"
	"sync"

	"gonum.org/v1/gonum/stat/distuv"

	"github.com/jana-veverkova/language-words-predictor-go/pkg/filehandler/csvx"
	"github.com/jana-veverkova/language-words-predictor-go/pkg/frequencydictionary"
	"github.com/pkg/errors"
)

const (
	// how many times a person has to see a word in order to remember it
	RequiredOccurence = 15
	// max number of words that will be tested
	WordLimit = 30000
)

// generates persons with different knowledge of words and saves them into files in targetDir
func CreateSourceFiles(targetDir string) {
	runtime.GOMAXPROCS(8)

	// get frequency dictionary
	freqDict, err := frequencydictionary.Get()
	if err != nil {
		fmt.Println(errors.WithStack(err))
	}

	// create sample of population and save
	nn := [][3]int{
		{1000, 1000000, 1000},
		{1000000, 3000000, 2000},
		{3000000, 6000000, 3000},
		{6000000, 10000000, 4000},
	}

	var wgGenerators sync.WaitGroup

	for i, n := range nn {
		wgGenerators.Add(1)
		go generatePopulationSample(targetDir, freqDict, i, n, &wgGenerators)
	}

	wgGenerators.Wait()
}

func generatePopulationSample(targetDir string, freqDict *frequencydictionary.FrequencyDictionary, ix int, n [3]int, wg *sync.WaitGroup) {
	defer wg.Done()

	var wgWriters sync.WaitGroup
	chRecords, chError := csvx.Write(fmt.Sprintf("%s/populationSample%d.csv", targetDir, ix), &wgWriters)

	var sample []string

	for i := n[0]; i < n[1]; i = i + n[2] {
		sample = generatePersonSample(freqDict, i)
		select {
		case chRecords <- sample:
			continue
		case err := <-chError:
			fmt.Println(errors.WithStack(err))
			close(chRecords)
			return
		}
	}

	close(chRecords)
	wgWriters.Wait()
}

func generatePersonSample(freqDict *frequencydictionary.FrequencyDictionary, wordsSeen int) []string {
	sample := make([]string, 0)

	var p float64
	var isKnown string

	for _, word := range freqDict.Words[:WordLimit] {
		p = getProbabilityOfWordOccurence(wordsSeen, word.Probability)
		isKnown = getBernoulliSample(p)
		sample = append(sample, isKnown)
	}

	return sample
}

func getProbabilityOfWordOccurence(wordsSeen int, wordProb float32) float64 {
	if wordsSeen < 1000 {
		binomial := distuv.Binomial{N: float64(wordsSeen), P: float64(wordProb)}
		return 1 - binomial.CDF(RequiredOccurence-1)
	} else {
		mean := float32(wordsSeen) * float32(wordProb)
		sd := math.Sqrt(float64(wordsSeen) * float64(wordProb) * (float64(1) - float64(wordProb)))
		z := (float64(RequiredOccurence-1) - 0.5 - float64(mean)) / sd
		normal := distuv.Normal{Mu: 0, Sigma: 1}
		return 1 - normal.CDF(z)
	}
}

func getBernoulliSample(p float64) string {
	bernoulli := distuv.Bernoulli{P: p}
	if s := bernoulli.Rand(); s == 1 {
		return "1"
	}
	return "0"
}
