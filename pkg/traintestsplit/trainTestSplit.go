package traintestsplit

import (
	"fmt"
	"math/rand"
	"os"
	"sync"

	"github.com/jana-veverkova/language-words-predictor-go/pkg/filehandler/csvx"
	"github.com/pkg/errors"
)

// goes through population sample and splits the sample 1:4 into train and test set

func Split(sourceDir string, targetDir string) {
	var wgWriters sync.WaitGroup

	chTest, chErrorTest := csvx.Write(targetDir+"/"+"test.csv", &wgWriters)
	chTrain, chErrorTrain := csvx.Write(targetDir+"/"+"train.csv", &wgWriters)

	items, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Println(errors.WithStack(err))
	}

L:
	for _, item := range items {
		var wgFileReaders sync.WaitGroup
		fmt.Printf("Reading %s/%s...", sourceDir, item.Name())

		wgFileReaders.Add(1)
		chRows, chError := csvx.Read(sourceDir+"/"+item.Name(), &wgFileReaders)
		for {
			select {
			case row, ok := <-chRows:
				if !ok {
					fmt.Print("...EOF \n")
					continue L
				}
				if rand.Intn(5) == 0 {
					select {
					case chTest <- row:
						continue
					case err := <-chErrorTest:
						fmt.Println(errors.WithStack(err))
						return
					}
				} else {
					select {
					case chTrain <- row:
						continue
					case err := <-chErrorTrain:
						fmt.Println(errors.WithStack(err))
						return
					}
				}
			case err := <-chError:
				fmt.Println(errors.WithStack(err))
				return
			}
		}
	}

	close(chTrain)
	close(chTest)
	wgWriters.Wait()
}
