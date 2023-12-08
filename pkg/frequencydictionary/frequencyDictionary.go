package frequencydictionary

import (
	"strconv"
	"sync"

	"github.com/jana-veverkova/language-words-predictor-go/pkg/filehandler/csvx"
)

type FrequencyDictionary struct {
	Words         []*WordOccurence
	TotalOccCount int
}

type WordOccurence struct {
	Word        string
	Occurence   int
	Probability float32
}

func Get() (*FrequencyDictionary, error) {
	sourceFile := "data/frequencyDictionaries/english.csv"

	dict := FrequencyDictionary{Words: []*WordOccurence{}, TotalOccCount: 0}

	var wgReaders sync.WaitGroup
	chRecords, chError := csvx.Read(sourceFile, &wgReaders)
L:
	for {
		select {
		case record, ok := <-chRecords:
			if !ok {
				break L
			}
			counts, err := strconv.Atoi(record[2])
			if err != nil {
				return &dict, err
			}
			dict.Words = append(dict.Words, &WordOccurence{Word: record[1], Occurence: counts})
			dict.TotalOccCount = dict.TotalOccCount + counts
		case err := <-chError:
			return &dict, err
		}
	}

	for _, word := range dict.Words {
		word.Probability = float32(word.Occurence) / float32(dict.TotalOccCount)
	}

	return &dict, nil
}
