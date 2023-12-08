package frequencydictionary

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/jana-veverkova/language-words-predictor-go/pkg/filehandler/csvx"
	"github.com/pkg/errors"
)

const wordPattern = "^[\\W|\\d]*([a-zA-Z]+(?:['’-][a-zA-Z]+)?)[\\W|\\d]*$"

var regex = regexp.MustCompile(wordPattern)

// takes raw files with sentences in target language and creates frequency dictionary file
func CreateFile(sourceDir string, targetFile string) {
	// creates frequency dictionary based on sentences in target language

	runtime.GOMAXPROCS(8)

	chLines := make(chan string, 1000)
	chWords1 := make(chan string, 1000)
	chWords2 := make(chan string, 1000)

	var wgFileReaders sync.WaitGroup
	var wgReceivers sync.WaitGroup
	var wgDictionaryMakers sync.WaitGroup

	items, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Println(errors.WithStack(err))
	}

	for _, item := range items {
		wgFileReaders.Add(1)
		go readFileLines(sourceDir+"/"+item.Name(), chLines, &wgFileReaders)
	}

	wgReceivers.Add(3)
	go extractWords(chLines, chWords1, chWords2, &wgReceivers)
	go extractWords(chLines, chWords1, chWords2, &wgReceivers)
	go extractWords(chLines, chWords1, chWords2, &wgReceivers)

	chDict1 := createDict(chWords1, &wgDictionaryMakers)
	chDict2 := createDict(chWords2, &wgDictionaryMakers)

	wgFileReaders.Wait()
	close(chLines)
	wgReceivers.Wait()
	close(chWords1)
	close(chWords2)

	dict1, dict2 := <-chDict1, <-chDict2
	dict := combine(dict1, dict2)

	wgDictionaryMakers.Wait()

	saveToFile(targetFile, dict)
}

func combine(dicts ...map[string]int) map[string]int {
	combinedDict := make(map[string]int)

	for _, dict := range dicts {
		for key, val := range dict {
			combinedDict[key] = val
		}
	}

	return combinedDict
}

func createDict(chWords chan string, wg *sync.WaitGroup) <-chan map[string]int {
	chDict := make(chan map[string]int)
	dict := make(map[string]int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for word := range chWords {
			dict[word] = dict[word] + 1
		}
		chDict <- dict
		close(chDict)
	}()

	return chDict
}

func readFileLines(fileName string, chLines chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(errors.WithStack(err))
	}

	defer file.Close()

	r := bufio.NewReader(file)

	for {
		line, _, err := r.ReadLine()
		if len(line) > 0 {
			chLines <- string(line)
		}
		if err != nil {
			break
		}
	}
}

func extractWords(chLines chan string, chWords1 chan string, chWords2 chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range chLines {
		words := strings.Split(line, " ")
		for _, word := range words {
			if w := matchWord(word); w != "" {
				c := strings.Compare(w, "l")
				if c < 0 {
					chWords1 <- w
				} else {
					chWords2 <- w
				}
			}
		}
	}
}

func matchWord(word string) string {
	matches := regex.FindStringSubmatch(word)
	if matches == nil {
		return ""
	}

	return strings.ToLower(matches[1])
}

func saveToFile(fileName string, dictionary map[string]int) {
	var wgWriter sync.WaitGroup
	chRecords, chError := csvx.Write(fileName, &wgWriter)

	orderedKeys := orderKeys(dictionary)

	for ix, word := range orderedKeys {
		select {
		case chRecords <- []string{fmt.Sprint(ix), word, fmt.Sprint(dictionary[word])}:
			continue
		case err := <-chError:
			fmt.Println(errors.WithStack(err))
			close(chRecords)
			wgWriter.Wait()
			return
		}
	}

	close(chRecords)
	wgWriter.Wait()
}

func orderKeys(dictionary map[string]int) []string {
	keys := make([]string, 0)

	for key := range dictionary {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return dictionary[keys[i]] > dictionary[keys[j]]
	})

	return keys
}
