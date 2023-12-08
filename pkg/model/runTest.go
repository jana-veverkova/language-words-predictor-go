package model

import (
	"fmt"

	"github.com/jana-veverkova/language-words-predictor-go/pkg/frequencydictionary"
	"github.com/jana-veverkova/language-words-predictor-go/pkg/model/prediction"
	"github.com/jana-veverkova/language-words-predictor-go/pkg/populationsample"
	"github.com/pkg/errors"
)

// asks words in console and predicts guess
func RunTest() {
	// get whole population sample
	trainSample, err := populationsample.GetAllSamples()
	if err != nil {
		fmt.Println(errors.WithStack(err))
		return
	}

	// get frequency dictionary
	freqDict, err := frequencydictionary.Get()
	if err != nil {
		fmt.Println(errors.WithStack(err))
		return
	}

	// verify word in terminal
	verify := func(wordRank int) bool {
		var response string
		word := freqDict.Words[wordRank].Word

		for {
			fmt.Printf("Word: %s. Do you know? y/n \n", word)

			_, err := fmt.Scanln(&response)
			if err != nil {
				fmt.Println(errors.WithStack(err))
			}

			if response == "y" {
				return true
			} else if response == "n" {
				return false
			} else {
				continue
			}
		}
	}

	setting := prediction.Setting{
		RangeOfWordsToAsk:  800,
		WordTestSize:       15,
		RequiredGuessRange: 300,
		MaxTestedWords:     100,
	}

	lowerGuess, upperGuess, wordsAsked, err := prediction.Predict(trainSample, verify, setting)
	if err != nil {
		fmt.Println(errors.WithStack(err))
		return
	}

	fmt.Printf("I asked %d words and my guess is: (%d, %d). \n", wordsAsked, lowerGuess, upperGuess)
}
