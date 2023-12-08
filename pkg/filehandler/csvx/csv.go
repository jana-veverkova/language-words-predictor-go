package csvx

import (
	"encoding/csv"
	"io"
	"os"
	"sync"
)

func Write(fileName string, wg *sync.WaitGroup) (chan<- []string, <-chan error) {
	chRecords := make(chan []string)
	chError := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		f, err := os.Create(fileName)
		if err != nil {
			chError <- err
			close(chError)
			return
		}

		defer f.Close()

		csvwriter := csv.NewWriter(f)

		for record := range chRecords {
			err := csvwriter.Write(record)
			if err != nil {
				chError <- err
				close(chError)
				return
			}
			csvwriter.Flush()
		}
	}()

	return chRecords, chError
}

func Read(fileName string, wg *sync.WaitGroup) (<-chan []string, <-chan error) {
	chRecords := make(chan []string)
	chError := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()

		f, err := os.Open(fileName)
		if err != nil {
			chError <- err
			close(chError)
			close(chRecords)
			return
		}

		defer f.Close()

		csvreader := csv.NewReader(f)

		for {
			record, err := csvreader.Read()
			if err != nil {
				if err != io.EOF {
					chError <- err
				}
				close(chRecords)
				return
			}

			chRecords <- record
		}
	}()

	return chRecords, chError
}
